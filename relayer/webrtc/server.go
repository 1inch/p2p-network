// Package webrtc starts webrtc server.
package webrtc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	pb "github.com/1inch/p2p-network/proto"
	"github.com/1inch/p2p-network/relayer/grpc"
	"github.com/pion/webrtc/v4"
	"google.golang.org/protobuf/proto"
)

const (
	maxRetries    = 5
	retryDelay    = 200 * time.Millisecond
	backoffFactor = 2
)

var (
	// ErrInvalidICEServer error represents invalid ICE server config.
	ErrInvalidICEServer = errors.New("invalid ICE server configuration")
	// ErrDataChannelNotFound error represents missing data channel.
	ErrDataChannelNotFound = errors.New("data channel not found for session")
	// ErrConnectionNotFound error represents missing connection.
	ErrConnectionNotFound = errors.New("connection not found for session")
)

// GRPCClient defines the interface for a gRPC client.
type GRPCClient interface {
	Execute(ctx context.Context, publicKey []byte, req *pb.ResolverRequest) (*pb.ResolverResponse, error)
	Close() error
}

// RegistryClient defines the interface for a node registry client.
type RegistryClient interface {
	GetResolver(publicKey []byte) (string, error)
}

// SDPRequest represents SDP request.
type SDPRequest struct {
	SessionID string
	Offer     webrtc.SessionDescription
	Response  chan *webrtc.SessionDescription
}

// ICECandidate represents ICECandidate request.
type ICECandidate struct {
	SessionID string
	Candidate webrtc.ICECandidate
}

// Server wraps the webrtc.Server.
type Server struct {
	cfg           Config
	logger        *slog.Logger
	ICEServer     string
	grpcClient    GRPCClient
	sdpRequests   <-chan SDPRequest
	iceCandidates <-chan ICECandidate
	connections   map[string]*webrtc.PeerConnection
	dataChannels  map[string]*webrtc.DataChannel
	mu            sync.RWMutex
}

// New initializes a new WebRTC server.
func New(
	cfg Config,
	logger *slog.Logger,
	iceServer string,
	client GRPCClient,
	sdpRequests <-chan SDPRequest,
	iceICECandidates <-chan ICECandidate,
) (*Server, error) {

	if iceServer == "" {
		return nil, ErrInvalidICEServer
	}

	return &Server{
		cfg:           cfg,
		sdpRequests:   sdpRequests,
		iceCandidates: iceICECandidates,
		ICEServer:     iceServer,
		grpcClient:    client,
		connections:   make(map[string]*webrtc.PeerConnection),
		dataChannels:  make(map[string]*webrtc.DataChannel),
		logger:        logger,
	}, nil
}

// HandleSDP processes an SDP offer, sets up a PeerConnection, and generates an SDP answer.
func (w *Server) HandleSDP(sessionID string, offer webrtc.SessionDescription) (*webrtc.SessionDescription, error) {
	w.logger.Debug("handle sdp", slog.String("sesionID", sessionID))

	pc, err := w.newPeerConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %w", err)
	}

	pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		w.logger.Debug("connection state change", slog.String("state", state.String()))

		if state == webrtc.PeerConnectionStateClosed || state == webrtc.PeerConnectionStateFailed {
			w.mu.Lock()
			delete(w.connections, sessionID)
			delete(w.dataChannels, sessionID)
			w.mu.Unlock()
		}
	})

	pc.OnICEConnectionStateChange(func(is webrtc.ICEConnectionState) {
		w.logger.Debug("ice connection state change", slog.String("sessionID", sessionID), slog.String("state", is.String()))
	})

	pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			w.logger.Debug("ice candidate gathering complete", slog.String("sessionID", sessionID))
			return
		}

		w.logger.Debug("ice candidate found", slog.String("sessionID", sessionID), slog.String("candidate", candidate.String()))
	})

	// Handle DataChannel setup.
	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		w.mu.Lock()
		w.dataChannels[sessionID] = dc
		w.mu.Unlock()

		dc.OnOpen(func() {
			w.logger.Debug("data channel opened", slog.String("sessionID", sessionID), slog.String("lable", dc.Label()))
		})

		w.handleDataChannel(dc)
	})

	// Set remote SDP description (offer).
	if err := pc.SetRemoteDescription(offer); err != nil {
		return nil, fmt.Errorf("failed to set remote description: %w", err)
	}

	// Generate and set local SDP description (answer).
	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create answer: %w", err)
	}

	gatherComplete := webrtc.GatheringCompletePromise(pc)

	if err := pc.SetLocalDescription(answer); err != nil {
		return nil, fmt.Errorf("failed to set local description: %w", err)
	}

	// Store the PeerConnection.
	w.mu.Lock()
	w.connections[sessionID] = pc
	w.mu.Unlock()

	<-gatherComplete

	return pc.LocalDescription(), nil
}

// Run starts the WebRTC server and processes SDP requests until context cancellation.
func (w *Server) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			w.cleanup()
			return nil
		case req, ok := <-w.sdpRequests:
			if !ok {
				return nil
			}
			answer, err := w.HandleSDP(req.SessionID, req.Offer)
			if err != nil {
				w.logger.Error("failed to process sdp offer", slog.Any("err", err))
				req.Response <- nil
			} else {
				req.Response <- answer
			}
		case req, ok := <-w.iceCandidates:
			if !ok {
				return nil
			}

			for attempt := 0; attempt < maxRetries; attempt++ {
				err := w.handleCandidate(req.SessionID, req.Candidate)
				if err == nil {
					break
				}

				if !errors.Is(err, ErrConnectionNotFound) {
					// Non-retryable errors
					w.logger.Error("failed to handle ICE candidate", slog.String("session_id", req.SessionID), slog.Any("err", err))
					break
				}

				w.logger.Warn("connection not found, retrying", slog.String("session_id", req.SessionID), slog.Int("attempt", attempt+1))

				select {
				case <-ctx.Done():
					w.logger.Warn("context cancelled during retry", slog.String("session_id", req.SessionID))
					return fmt.Errorf("context cancelled: %w", err)
				case <-time.After(retryWithBackoff(attempt)):
				}
			}
		}
	}
}

// SendMessage sends a message over the DataChannel associated with the given session ID.
func (w *Server) SendMessage(sessionID, message string) error {
	w.mu.RLock()
	dc, ok := w.dataChannels[sessionID]
	w.mu.RUnlock()

	if !ok {
		return ErrDataChannelNotFound
	}

	if err := dc.SendText(message); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// GetConnection retrieves a PeerConnection by session ID.
func (w *Server) GetConnection(sessionID string) (*webrtc.PeerConnection, bool) {
	w.mu.RLock()
	pc, ok := w.connections[sessionID]
	w.mu.RUnlock()
	return pc, ok
}

// GetAllConnections retrieves all active PeerConnection session IDs.
func (w *Server) GetAllConnections() []string {
	w.mu.RLock()
	sessions := make([]string, 0, len(w.connections))
	for sessionID := range w.connections {
		sessions = append(sessions, sessionID)
	}
	w.mu.RUnlock()
	return sessions
}

func (w *Server) handleDataChannel(dc *webrtc.DataChannel) {
	dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		var message pb.IncomingMessage
		if err := proto.Unmarshal(msg.Data, &message); err != nil {
			respMessage := &pb.OutgoingMessage{
				Result: &pb.OutgoingMessage_Error{
					Error: &pb.Error{
						Code:    pb.ErrorCode_ERR_INVALID_MESSAGE_FORMAT,
						Message: fmt.Sprintf("failed to unmarshal protobuf message: %v", err),
					},
				},
			}
			w.logger.Error("failed to unmarshal protobuf message", slog.Any("err", err))
			if sendErr := w.sendResponse(dc, respMessage); sendErr != nil {
				w.logger.Error("failed to send invalid message error response", slog.Any("err", sendErr))
			}
			return
		}

		w.logger.Debug("received message", slog.Any("request", message.Request), slog.String("publicKeys", fmt.Sprintf("%x", message.PublicKeys)))

		respMessageChan := make(chan *pb.OutgoingMessage, len(message.PublicKeys))

		waitGroupForRequestGoroutine := &sync.WaitGroup{}
		for _, publicKey := range message.PublicKeys {
			waitGroupForRequestGoroutine.Add(1)
			go func() {
				w.retryGetResponseFromResolver(publicKey, message.Request, respMessageChan)
				defer waitGroupForRequestGoroutine.Done()
			}()
		}
		// Waiting when all requests return response (correct, with error, etc)
		waitGroupForRequestGoroutine.Wait()

		respMessage := w.tryFindCorrectResponse(respMessageChan)

		if err := w.sendResponse(dc, respMessage); err != nil {
			respMessage := &pb.OutgoingMessage{
				Result: &pb.OutgoingMessage_Error{
					Error: &pb.Error{
						Code:    pb.ErrorCode_ERR_DATA_CHANNEL_SEND_FAILED,
						Message: fmt.Sprintf("Failed to send response: %v", err),
					},
				},
				PublicKey: respMessage.PublicKey,
			}
			w.logger.Error("failed to send response", slog.Any("err", err))
			if sendErr := w.sendResponse(dc, respMessage); sendErr != nil {
				w.logger.Error("failed to send data channel error response", slog.Any("err", sendErr))
			}
		}
	})
}

func (w *Server) sendResponse(dc *webrtc.DataChannel, message *pb.OutgoingMessage) error {
	respBytes, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal protobuf response: %w", err)
	}

	if err := dc.Send(respBytes); err != nil {
		return fmt.Errorf("failed to send response: %w", err)
	}

	return nil
}

func (w *Server) handleCandidate(sessionID string, candidate webrtc.ICECandidate) error {
	w.logger.Debug("handled ice candidate", slog.String("sessionID", sessionID), slog.String("candidate", candidate.String()))

	w.mu.RLock()
	conn, ok := w.connections[sessionID]
	w.mu.RUnlock()

	if !ok {
		return fmt.Errorf("%w: session_id=%s", ErrConnectionNotFound, sessionID)
	}

	err := conn.AddICECandidate(candidate.ToJSON())
	if err != nil {
		return fmt.Errorf("failed to add ICE candidate: %w", err)
	}

	return nil
}

func retryWithBackoff(attempt int) time.Duration {
	return retryDelay * time.Duration(backoffFactor^attempt)
}

func (w *Server) cleanup() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for sessionID, pc := range w.connections {
		if err := pc.Close(); err != nil {
			w.logger.Error("failed to close peer connection", slog.String("session_id", sessionID), slog.Any("err", err))
		}
	}
	w.connections = make(map[string]*webrtc.PeerConnection)
	w.dataChannels = make(map[string]*webrtc.DataChannel)

	if err := w.grpcClient.Close(); err != nil {
		w.logger.Error("failed to close gRPC client", slog.Any("err", err))
	}
}

func mapErrorToCodeAndMessage(err error) (pb.ErrorCode, string) {
	if errors.Is(err, grpc.ErrResolverLookupFailed) {
		return pb.ErrorCode_ERR_RESOLVER_LOOKUP_FAILED, fmt.Sprintf("resolver lookup failed: %s", err.Error())
	} else if errors.Is(err, grpc.ErrGRPCExecutionFailed) {
		return pb.ErrorCode_ERR_GRPC_EXECUTION_FAILED, fmt.Sprintf("grpc execution failed: %s", err.Error())
	}
	return pb.ErrorCode_ERR_GRPC_EXECUTION_FAILED, fmt.Sprintf("unexpected error: %s", err.Error())
}

func (w *Server) retryGetResponseFromResolver(publicKey []byte, request *pb.ResolverRequest, respChan chan *pb.OutgoingMessage) {
	w.logger.Debug("start request to resolver", slog.Any("public_key", string(publicKey)))

	resp := &pb.OutgoingMessage{}
	retryRequest := w.cfg.RetryConfig.Count
	requestSleepInterval := w.cfg.RetryConfig.Interval

	if !w.cfg.RetryConfig.Enabled {
		w.logger.Debug("retry requests is disabled, set parameters for one")
		retryRequest = 1
		requestSleepInterval = time.Duration(0)
	}
	for attempt := range retryRequest {
		w.logger.Debug("try get response from resolver", slog.Any("attempt", attempt+1), slog.Any("publicKey", string(publicKey)))
		resp = w.tryGetResponseFromResolver(publicKey, request)

		// if response without error, there is break cycle for attempts and put the message in chan
		if resp.GetError() == nil {
			w.logger.Info("success response from resolver")
			w.logger.Debug("response from resolver without error, put it in channel", slog.Any("publicKey", string(publicKey)))
			break
		}

		// But if response with error, check which error that, start sleep and continue cycle for attempts
		switch {
		case w.isErrorForRetry(resp.GetError().Code):
			{
				w.logger.Debug("resolver returned supported error for retry")
				w.logger.Debug("start sleep for interval", slog.Any("time_duration", requestSleepInterval))
				time.Sleep(requestSleepInterval)
			}
		// there is no point retry get some response from resolvers, because nothing will change
		// for example, error when dapps send request with not supported method in api handler
		default:
			{
				w.logger.Error("unsupported error from resolver, retry is not possible", slog.Any("err", resp.GetError().Message))
				break
			}
		}
	}
	w.logger.Debug("put OutgoingMessage to response channel", slog.Any("publicKey", resp.PublicKey))
	respChan <- resp
}

func (w *Server) isErrorForRetry(errorCode pb.ErrorCode) bool {
	return errorCode == pb.ErrorCode_ERR_DATA_CHANNEL_SEND_FAILED ||
		errorCode == pb.ErrorCode_ERR_GRPC_EXECUTION_FAILED ||
		errorCode == pb.ErrorCode_ERR_RESOLVER_LOOKUP_FAILED
}

func (w *Server) tryGetResponseFromResolver(publicKey []byte, request *pb.ResolverRequest) *pb.OutgoingMessage {
	var respMessage *pb.OutgoingMessage
	response, err := w.grpcClient.Execute(context.Background(), publicKey, request)
	if err != nil {
		var errorCode, errorMsg = mapErrorToCodeAndMessage(err)

		respMessage = &pb.OutgoingMessage{
			Result: &pb.OutgoingMessage_Error{
				Error: &pb.Error{
					Code:    errorCode,
					Message: errorMsg,
				},
			},
			PublicKey: publicKey,
		}
		w.logger.Error("failed to execute gRPC request", slog.Any("err", err))
	} else {
		respMessage = &pb.OutgoingMessage{
			Result: &pb.OutgoingMessage_Response{
				Response: response,
			},
			PublicKey: publicKey,
		}
	}

	return respMessage
}

// try find correct response from resolvers, if cant find correct - return someone, probably with error
func (w *Server) tryFindCorrectResponse(chanWithResp chan *pb.OutgoingMessage) *pb.OutgoingMessage {
	resps := make([]*pb.OutgoingMessage, len(chanWithResp))

	// rewrite in array responses from channel and check this response is success
	for i := range resps {
		resps[i] = <-chanWithResp

		w.logger.Debug("check is response successful", slog.Any("publicKey", resps[i].PublicKey))
		if resps[i].GetError() == nil {
			return resps[i]
		}
	}

	// if successful response not found, return first respons with error
	return resps[0]
}

// create and configure new peer connection
func (w *Server) newPeerConnection() (*webrtc.PeerConnection, error) {
	s := webrtc.SettingEngine{}

	if w.cfg.PeerPortConfig.Enabled {
		err := s.SetEphemeralUDPPortRange(w.cfg.PeerPortConfig.Min, w.cfg.PeerPortConfig.Max)
		if err != nil {
			w.logger.Error("failed set peers port range", slog.Any("err", err.Error()))
			return nil, err
		}
	}

	api := webrtc.NewAPI(webrtc.WithSettingEngine(s))

	return api.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{{URLs: []string{w.ICEServer}}},
	})
}

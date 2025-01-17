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

	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{{URLs: []string{w.ICEServer}}},
	})
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
			if sendErr := w.sendResponse(dc, respMessage); sendErr != nil {
				w.logger.Error("failed to send invalid message error response", slog.Any("err", sendErr))
			}
			return
		}

		var wg sync.WaitGroup
		for _, publicKey := range message.PublicKeys {
			wg.Add(1)
			go func(publicKey []byte) {
				defer wg.Done()

				response, err := w.grpcClient.Execute(context.Background(), publicKey, message.Request)
				if err != nil {
					var errorCode, errorMsg = mapErrorToCodeAndMessage(err)

					respMessage := &pb.OutgoingMessage{
						Result: &pb.OutgoingMessage_Error{
							Error: &pb.Error{
								Code:    errorCode,
								Message: errorMsg,
							},
						},
						PublicKey: publicKey,
					}
					if sendErr := w.sendResponse(dc, respMessage); sendErr != nil {
						w.logger.Error("failed to send error response", slog.Any("err", sendErr))
					}
					return
				}

				respMessage := &pb.OutgoingMessage{
					Result: &pb.OutgoingMessage_Response{
						Response: response,
					},
					PublicKey: publicKey,
				}

				if err := w.sendResponse(dc, respMessage); err != nil {
					respMessage := &pb.OutgoingMessage{
						Result: &pb.OutgoingMessage_Error{
							Error: &pb.Error{
								Code:    pb.ErrorCode_ERR_DATA_CHANNEL_SEND_FAILED,
								Message: fmt.Sprintf("Failed to send response: %v", err),
							},
						},
						PublicKey: publicKey,
					}
					if sendErr := w.sendResponse(dc, respMessage); sendErr != nil {
						w.logger.Error("failed to send data channel error response", slog.Any("err", sendErr))
					}
				}
			}(publicKey)
		}
		wg.Wait()
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

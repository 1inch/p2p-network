// Package webrtc starts webrtc server.
package webrtc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	pbrelayer "github.com/1inch/p2p-network/proto/relayer"
	pbresolver "github.com/1inch/p2p-network/proto/resolver"
	"github.com/1inch/p2p-network/relayer/metrics"
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

// Option represents configuration of some server parameters
type Option func(*Server)

// GRPCClient defines the interface for a gRPC client.
type GRPCClient interface {
	Execute(ctx context.Context, publicKey []byte, req *pbresolver.ResolverRequest) (*pbresolver.ResolverResponse, error)
	Close() error
}

// RegistryClient defines the interface for a node registry client.
type RegistryClient interface {
	GetResolver(publicKey []byte) (string, error)
}

// SDPRequest represents SDP request.
type SDPRequest struct {
	SessionID    string
	Offer        webrtc.SessionDescription
	CandidateURL string
	Response     chan *webrtc.SessionDescription
}

// ICECandidate represents ICECandidate request.
type ICECandidate struct {
	SessionID string              `json:"session_id"`
	Candidate webrtc.ICECandidate `json:"candidate"`
}

// Server wraps the webrtc.Server.
type Server struct {
	useTrickleICE bool
	retryOpt      *Retry
	peerPortOpt   *PeerRangePort
	logger        *slog.Logger
	iceServers    []webrtc.ICEServer
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
	iceServers []webrtc.ICEServer,
	client GRPCClient,
	sdpRequests <-chan SDPRequest,
	iceICECandidates <-chan ICECandidate,
	options ...Option,
) (*Server, error) {

	if iceServers == nil {
		return nil, ErrInvalidICEServer
	}

	srv := &Server{
		sdpRequests:   sdpRequests,
		iceCandidates: iceICECandidates,
		iceServers:    iceServers,
		grpcClient:    client,
		connections:   make(map[string]*webrtc.PeerConnection),
		dataChannels:  make(map[string]*webrtc.DataChannel),
		logger:        logger,
	}

	for _, opt := range options {
		opt(srv)
	}

	return srv, nil
}

// WithRetry added retry option in webrtc server
func WithRetry(option Retry) Option {
	return func(s *Server) {
		s.retryOpt = &option
	}
}

// WithPeerPort added peers port option in webrtc server
func WithPeerPort(option PeerRangePort) Option {
	return func(s *Server) {
		s.peerPortOpt = &option
	}
}

// WithTrickleICE added send request about candidate
func WithTrickleICE() Option {
	return func(s *Server) {
		s.useTrickleICE = true
	}
}

// HandleSDP processes an SDP offer, sets up a PeerConnection, and generates an SDP answer.
func (w *Server) HandleSDP(candidateURL, sessionID string, offer webrtc.SessionDescription) (*webrtc.SessionDescription, error) {
	start := time.Now()

	w.logger.Debug("handle sdp", slog.String("sesionID", sessionID))

	pc, err := w.newPeerConnection()
	if err != nil {
		metrics.SdpNegotiationTotal.WithLabelValues("failure").Inc()
		return nil, fmt.Errorf("failed to create peer connection: %w", err)
	}

	pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		w.logger.Debug("connection state change", slog.String("state", state.String()))

		if state == webrtc.PeerConnectionStateClosed || state == webrtc.PeerConnectionStateFailed {
			w.mu.Lock()
			delete(w.connections, sessionID)
			delete(w.dataChannels, sessionID)
			w.mu.Unlock()
			metrics.ActivePeerConnections.Dec()
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
		if w.useTrickleICE {
			go w.sendCandidate(candidateURL, sessionID, *candidate)
		}
	})

	// Handle DataChannel setup.
	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		w.mu.Lock()
		w.dataChannels[sessionID] = dc
		w.mu.Unlock()

		dc.OnOpen(func() {
			w.logger.Debug("data channel opened", slog.String("sessionID", sessionID), slog.String("lable", dc.Label()))
		})

		w.handleDataChannel(dc, sessionID)
	})

	// Set remote SDP description (offer).
	if err := pc.SetRemoteDescription(offer); err != nil {
		metrics.SdpNegotiationTotal.WithLabelValues("failure").Inc()
		return nil, fmt.Errorf("failed to set remote description: %w", err)
	}

	// Generate and set local SDP description (answer).
	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		metrics.SdpNegotiationTotal.WithLabelValues("failure").Inc()
		return nil, fmt.Errorf("failed to create answer: %w", err)
	}

	gatherComplete := webrtc.GatheringCompletePromise(pc)

	if err := pc.SetLocalDescription(answer); err != nil {
		metrics.SdpNegotiationTotal.WithLabelValues("failure").Inc()
		return nil, fmt.Errorf("failed to set local description: %w", err)
	}

	// Store the PeerConnection.
	w.mu.Lock()
	w.connections[sessionID] = pc
	w.mu.Unlock()

	metrics.ActivePeerConnections.Inc()

	if !w.useTrickleICE {
		<-gatherComplete
	}

	metrics.SdpNegotiationTotal.WithLabelValues("success").Inc()
	metrics.SdpNegotiationDuration.Observe(time.Since(start).Seconds())

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
			answer, err := w.HandleSDP(req.CandidateURL, req.SessionID, req.Offer)
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

func (w *Server) handleDataChannel(dc *webrtc.DataChannel, sessionID string) {
	dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		start := time.Now()
		metrics.DataChannelMessagesReceived.WithLabelValues(sessionID).Inc()

		var message pbrelayer.IncomingMessage
		if err := proto.Unmarshal(msg.Data, &message); err != nil {
			respMessage := w.buildOutgoingMessageWithErr(
				[]byte{},
				pbrelayer.ErrorCode_ERR_INVALID_MESSAGE_FORMAT,
				fmt.Sprintf("failed to unmarshal protobuf message: %v", err))

			w.logger.Error("failed to unmarshal protobuf message", slog.Any("err", err))

			if sendErr := w.sendResponse(dc, respMessage); sendErr != nil {
				w.logger.Error("failed to send invalid message error response", slog.Any("err", sendErr))
			}
			metrics.DataChannelMessagesSent.WithLabelValues(sessionID, "failed").Inc()

			return
		}

		w.logger.Debug("received message", slog.Any("request", message.Request), slog.String("publicKeys", fmt.Sprintf("%x", message.PublicKeys)))

		doneChan := make(chan bool)
		respChan := make(chan *pbrelayer.OutgoingMessage)

		for _, publicKey := range message.PublicKeys {
			go w.retryGetResponseFromResolver(publicKey, message.Request, doneChan, respChan)
		}

		respMessage := <-respChan
		if err := w.sendResponse(dc, respMessage); err != nil {
			w.logger.Error("failed to send response", slog.Any("err", err))
		}
		status := "success"
		if respMessage.GetError() != nil {
			status = "failed"
			w.logger.Error("failed to send response", slog.Any("err", respMessage.GetError().Message))
		}
		metrics.DataChannelMessagesSent.WithLabelValues(sessionID, status).Inc()

		latency := time.Since(start).Seconds()
		metrics.DataChannelLatency.WithLabelValues(sessionID).Observe(latency)
	})
}

func (w *Server) sendResponse(dc *webrtc.DataChannel, message *pbrelayer.OutgoingMessage) error {
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

func (w *Server) retryGetResponseFromResolver(publicKey []byte, request *pbresolver.ResolverRequest, doneChan chan bool, respChan chan *pbrelayer.OutgoingMessage) {
	w.logger.Debug("start request to resolver", slog.Any("public_key", string(publicKey)))

	resp := &pbrelayer.OutgoingMessage{}
	retryRequest := uint8(1)
	requestSleepInterval := time.Duration(0)

	if w.retryOpt != nil {
		w.logger.Debug("retry requests is enabled, set parameters")
		retryRequest = w.retryOpt.Count
		requestSleepInterval = w.retryOpt.Interval
	}
	for attempt := range retryRequest {
		// check doneChan before start try get response from resolver
		select {
		case <-doneChan:
			{
				w.logger.Debug("some gorotine returned response, stop this gorotine")
				return
			}
		default:
			{
				w.logger.Debug("try get response from resolver", slog.Any("attempt", attempt+1), slog.Any("publicKey", fmt.Sprintf("%x", publicKey)))
				resolverResponse, err := w.grpcClient.Execute(context.Background(), publicKey, request)

				// if grpc call return error, try retry
				if err != nil {
					// put OutgoingMessage with error
					resp = &pbrelayer.OutgoingMessage{
						PublicKey: publicKey,
						Result: &pbrelayer.OutgoingMessage_Error{
							Error: &pbrelayer.Error{
								Code:    pbrelayer.ErrorCode_ERR_GRPC_EXECUTION_FAILED,
								Message: fmt.Sprintf("failed call execute: %v", err),
							},
						},
					}

					w.logger.Debug("resolver returned supported error for retry")
					w.logger.Debug("start sleep for interval", slog.Any("time_duration", requestSleepInterval))
					time.Sleep(requestSleepInterval)
					continue
				}
				// At this moment, if resolver return any error in ResolverResponse.Error that error not for retry
				resp = &pbrelayer.OutgoingMessage{
					PublicKey: publicKey,
					Result: &pbrelayer.OutgoingMessage_Response{
						Response: resolverResponse,
					},
				}
				break
			}
		}
	}
	w.logger.Debug("put OutgoingMessage to response channel", slog.Any("publicKey", fmt.Sprintf("%x", publicKey)))
	respChan <- resp
	doneChan <- true
}

// create and configure new peer connection
func (w *Server) newPeerConnection() (*webrtc.PeerConnection, error) {
	s := webrtc.SettingEngine{}

	if w.peerPortOpt != nil {
		err := s.SetEphemeralUDPPortRange(w.peerPortOpt.Min, w.peerPortOpt.Max)
		if err != nil {
			w.logger.Error("failed set peers port range", slog.Any("err", err.Error()))
			return nil, err
		}
	}

	api := webrtc.NewAPI(webrtc.WithSettingEngine(s))

	return api.NewPeerConnection(webrtc.Configuration{
		ICEServers: w.iceServers,
	})
}

func (w *Server) sendCandidate(candidateURL string, sessionID string, candidate webrtc.ICECandidate) {
	start := time.Now()
	payload := ICECandidate{
		SessionID: sessionID,
		Candidate: candidate,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		w.logger.Error("failed to marshal candidate payload", slog.String("sessionID", sessionID), slog.Any("err", err))
		return
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", candidateURL, bytes.NewReader(data))
	if err != nil {
		w.logger.Error("failed to create candidate request", slog.String("sessionID", sessionID), slog.Any("err", err))
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.logger.Error("failed to send candidate", slog.String("sessionID", sessionID), slog.Any("err", err))
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			w.logger.Error("failed to close response body", slog.String("sessionID", sessionID), slog.Any("err", err))
		}
	}()

	duration := time.Since(start).Seconds()
	metrics.IceCandidateSendDuration.WithLabelValues(sessionID).Observe(duration)
	metrics.IceCandidateSentTotal.WithLabelValues(sessionID, resp.Status).Inc()

	w.logger.Debug("candidate sent", slog.String("sessionID", sessionID), slog.Any("candidate", candidate), slog.String("status", resp.Status))
}

func (w *Server) buildOutgoingMessageWithErr(publickKey []byte, errCode pbrelayer.ErrorCode, errMsg string) *pbrelayer.OutgoingMessage {
	return &pbrelayer.OutgoingMessage{
		PublicKey: publickKey,
		Result: &pbrelayer.OutgoingMessage_Error{
			Error: &pbrelayer.Error{
				Code:    errCode,
				Message: errMsg,
			},
		},
	}
}

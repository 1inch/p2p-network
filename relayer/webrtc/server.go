// Package webrtc starts webrtc server.
package webrtc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/pion/webrtc/v3"
)

var (
	// ErrInvalidICEServer error represents invalid ICE server config.
	ErrInvalidICEServer = errors.New("invalid ICE server configuration")
)

type SDPRequest struct {
	SessionID string
	Offer     webrtc.SessionDescription
	Response  chan *webrtc.SessionDescription
}

type Server struct {
	logger       *slog.Logger
	ICEServer    string
	sdpRequests  <-chan SDPRequest
	connections  map[string]*webrtc.PeerConnection
	dataChannels map[string]*webrtc.DataChannel
	mu           sync.Mutex
	OnMessage    func(sessionID, message string) // Callback for incoming messages
}

// New initializes a new WebRTC server.
func New(logger *slog.Logger, iceServer string, sdpRequests <-chan SDPRequest) (*Server, error) {
	if iceServer == "" {
		return nil, fmt.Errorf("invalid ICE server configuration")
	}

	return &Server{
		sdpRequests:  sdpRequests,
		ICEServer:    iceServer,
		connections:  make(map[string]*webrtc.PeerConnection),
		dataChannels: make(map[string]*webrtc.DataChannel),
		logger:       logger,
	}, nil
}

// HandleSDP processes an SDP offer, sets up a PeerConnection, and generates an SDP answer.
func (w *Server) HandleSDP(sessionID string, offer webrtc.SessionDescription) (*webrtc.SessionDescription, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{{URLs: []string{w.ICEServer}}},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %w", err)
	}

	// Handle DataChannel setup
	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		w.logger.Info("data channel opened", slog.String("label", dc.Label()))

		w.mu.Lock()
		w.dataChannels[sessionID] = dc
		w.mu.Unlock()

		// Handle incoming messages
		dc.OnMessage(func(msg webrtc.DataChannelMessage) {
			message := string(msg.Data)
			w.logger.Debug("received message", slog.String("session_id", sessionID), slog.String("message", message))

			if w.OnMessage != nil {
				w.OnMessage(sessionID, message)
			}
		})
	})

	o, err := pc.CreateOffer(nil)
	if err != nil {
		return nil, fmt.Errorf("offer create offer error: %w", err)
	}

	// Set remote SDP description (offer)
	if err := pc.SetRemoteDescription(o); err != nil {
		return nil, fmt.Errorf("failed to set remote description: %w", err)
	}

	// Generate and set local SDP description (answer)
	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create answer: %w", err)
	}
	if err := pc.SetLocalDescription(answer); err != nil {
		return nil, fmt.Errorf("failed to set local description: %w", err)
	}

	// Store the PeerConnection
	w.connections[sessionID] = pc
	w.logger.Info("peer connection established", slog.String("session_id", sessionID))
	return &answer, nil
}

// Run starts the WebRTC server and processes SDP requests until context cancellation.
func (w *Server) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			w.logger.Info("shutting down webrtc server")
			w.cleanup()
			return nil
		case req := <-w.sdpRequests:
			answer, err := w.HandleSDP(req.SessionID, req.Offer)
			if err != nil {
				w.logger.Error("Failed to process sdp offer", slog.String("error", err.Error()))
				req.Response <- nil
			} else {
				req.Response <- answer
			}
		}
	}
}

// SendMessage sends a message over the DataChannel associated with the given session ID.
func (w *Server) SendMessage(sessionID, message string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	dc, ok := w.dataChannels[sessionID]
	if !ok {
		return fmt.Errorf("data channel not found for session: %s", sessionID)
	}

	if err := dc.SendText(message); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	w.logger.Debug("message sent", slog.String("session_id", sessionID), slog.String("message", message))
	return nil
}

// GetConnection retrieves a PeerConnection by session ID.
func (w *Server) GetConnection(sessionID string) (*webrtc.PeerConnection, bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	pc, ok := w.connections[sessionID]
	return pc, ok
}

// GetAllConnections retrieves all active PeerConnection session IDs.
func (w *Server) GetAllConnections() []string {
	w.mu.Lock()
	defer w.mu.Unlock()

	sessions := make([]string, 0, len(w.connections))
	for sessionID := range w.connections {
		sessions = append(sessions, sessionID)
	}
	return sessions
}

func (w *Server) cleanup() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for sessionID, pc := range w.connections {
		w.logger.Info("closing peer connection", slog.String("session_id", sessionID))
		if err := pc.Close(); err != nil {
			w.logger.Error("failed to close peer connection", slog.String("session_id", sessionID), slog.String("error", err.Error()))
		}
		delete(w.connections, sessionID)
	}

	// Clear dataChannels
	w.dataChannels = make(map[string]*webrtc.DataChannel)
	w.logger.Info("webrtc server shutdown completed")
}

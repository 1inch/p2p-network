// Package webrtc starts webrtc server.
package webrtc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/pion/webrtc/v4"
)

var (
	// ErrInvalidICEServer error represents invalid ICE server config.
	ErrInvalidICEServer = errors.New("invalid ICE server configuration")
	// ErrDataChannelNotFound error represents missing data channel.
	ErrDataChannelNotFound = errors.New("data channel not found for session")
)

// SDPRequest represents SDP request.
type SDPRequest struct {
	SessionID string
	Offer     webrtc.SessionDescription
	Response  chan *webrtc.SessionDescription
}

// Server wraps the webrtc.Server.
type Server struct {
	logger       *slog.Logger
	ICEServer    string
	sdpRequests  <-chan SDPRequest
	connections  map[string]*webrtc.PeerConnection
	dataChannels map[string]*webrtc.DataChannel
	mu           sync.RWMutex
	OnMessage    func(sessionID, message string) // Callback for incoming messages
}

// New initializes a new WebRTC server.
func New(logger *slog.Logger, iceServer string, sdpRequests <-chan SDPRequest) (*Server, error) {
	if iceServer == "" {
		return nil, ErrInvalidICEServer
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

	// Handle DataChannel setup.
	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		w.logger.Info("establishing data channel", slog.String("label", dc.Label()))

		w.mu.Lock()
		w.dataChannels[sessionID] = dc
		w.mu.Unlock()

		dc.OnOpen(func() {
			w.logger.Info("data channel opened", slog.String("label", dc.Label()))
		})

		// Handle incoming messages.
		dc.OnMessage(func(msg webrtc.DataChannelMessage) {
			message := string(msg.Data)
			w.logger.Debug("received message", slog.String("session_id", sessionID), slog.String("message", message))

			if w.OnMessage != nil {
				w.OnMessage(sessionID, message)
			}
		})
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
	if err := pc.SetLocalDescription(answer); err != nil {
		return nil, fmt.Errorf("failed to set local description: %w", err)
	}

	// Store the PeerConnection.
	w.mu.Lock()
	w.connections[sessionID] = pc
	w.mu.Unlock()

	return &answer, nil
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
				w.logger.Error("Failed to process SDP offer", slog.String("error", err.Error()))
				req.Response <- nil
			} else {
				req.Response <- answer
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

func (w *Server) cleanup() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for sessionID, pc := range w.connections {
		if err := pc.Close(); err != nil {
			w.logger.Error("failed to close peer connection", slog.String("session_id", sessionID), slog.String("error", err.Error()))
		}
	}
	w.connections = make(map[string]*webrtc.PeerConnection)
	w.dataChannels = make(map[string]*webrtc.DataChannel)
}

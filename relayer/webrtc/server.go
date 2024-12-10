// Package webrtc starts webrtc server.
package webrtc

import (
	"errors"
)

var (
	// ErrInvalidICEServer error represents invalid ICE server config.
	ErrInvalidICEServer = errors.New("invalid ICE server configuration")
)

// Server represents the WebRTC server.
type Server struct {
	ICEServer string
}

// New initializes a new WebRTC server.
func New(iceServer string) (*Server, error) {
	if iceServer == "" {
		return nil, ErrInvalidICEServer
	}

	return &Server{ICEServer: iceServer}, nil
}

// Run starts the WebRTC server.
func (w *Server) Run() error {
	return nil
}

package webrtc_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/pion/webrtc/v4"
	"github.com/stretchr/testify/assert"

	relayerwebrtc "github.com/1inch/p2p-network/relayer/webrtc"
)

func TestWebRTCServer_HandleSDP(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	sdpRequests := make(chan relayerwebrtc.SDPRequest, 1)

	server, err := relayerwebrtc.New(logger, "stun:stun.l.google.com:19302", nil, sdpRequests)
	assert.NoError(t, err, "Failed to create WebRTC server")

	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{{URLs: []string{
			"stun:stun.l.google.com:19302",
		}}},
	})
	assert.NoError(t, err, "Failed to create dummy PeerConnection")

	_, err = peerConnection.CreateDataChannel("data", nil)
	assert.NoError(t, err, "Failed to create dummy DataChannel")

	offer, err := peerConnection.CreateOffer(nil)
	assert.NoError(t, err, "Failed to create SDP offer")

	err = peerConnection.SetLocalDescription(offer)
	assert.NoError(t, err, "Failed to set local description for dummy PeerConnection")

	// Simulate SDP request.
	responseChan := make(chan *webrtc.SessionDescription)
	sdpRequests <- relayerwebrtc.SDPRequest{
		SessionID: "test-session",
		Offer:     offer,
		Response:  responseChan,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := server.Run(ctx)
		assert.NoError(t, err, "WebRTC server exited with error")
	}()

	// Wait for the SDP answer.
	answer := <-responseChan
	assert.NotNil(t, answer, "Expected SDP answer to be non-nil")

	// Validate that the PeerConnection was created.
	pc, ok := server.GetConnection("test-session")
	assert.True(t, ok, "Expected PeerConnection to be stored in server")
	state := pc.ConnectionState()
	assert.Condition(t, func() bool {
		return state == webrtc.PeerConnectionStateNew || state == webrtc.PeerConnectionStateConnecting
	}, "Expected PeerConnection state to be 'new' or 'connecting', got: %s", state)
}

func TestWebRTCServer_Run_CleanupOnContextCancel(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	sdpRequests := make(chan relayerwebrtc.SDPRequest, 1)

	server, err := relayerwebrtc.New(logger, "stun:stun.l.google.com:19302", nil, sdpRequests)
	assert.NoError(t, err, "Failed to create WebRTC server")

	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	assert.NoError(t, err, "Failed to create dummy PeerConnection")

	_, err = peerConnection.CreateDataChannel("data", nil)
	assert.NoError(t, err, "Failed to create dummy DataChannel")

	offer, err := peerConnection.CreateOffer(nil)
	assert.NoError(t, err, "Failed to create SDP offer")

	err = peerConnection.SetLocalDescription(offer)
	assert.NoError(t, err, "Failed to set local description for dummy PeerConnection")

	// Simulate SDP request.
	responseChan := make(chan *webrtc.SessionDescription)
	sdpRequests <- relayerwebrtc.SDPRequest{
		SessionID: "test-session",
		Offer:     offer,
		Response:  responseChan,
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err := server.Run(ctx)
		assert.NoError(t, err, "WebRTC server exited with error")
	}()

	// Wait for the SDP answer.
	answer := <-responseChan
	assert.NotNil(t, answer, "Expected SDP answer to be non-nil")

	// Validate that the PeerConnection was created.
	_, ok := server.GetConnection("test-session")
	assert.True(t, ok, "Expected PeerConnection to be stored in server")

	// Cancel the context to trigger cleanup.
	cancel()

	// Allow some time for cleanup to complete.
	time.Sleep(100 * time.Millisecond)

	// Validate that connections are cleaned up
	_, exists := server.GetConnection("test-session")
	assert.False(t, exists, "Expected PeerConnection to be removed after cleanup")
}

func TestWebRTCServer_Run_Shutdown(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	sdpRequests := make(chan relayerwebrtc.SDPRequest, 1)

	server, err := relayerwebrtc.New(logger, "stun:stun.l.google.com:19302", nil, sdpRequests)
	assert.NoError(t, err, "Failed to create WebRTC server")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := server.Run(ctx)
		assert.NoError(t, err, "WebRTC server exited with error")
	}()

	// Simulate server shutdown.
	time.Sleep(100 * time.Millisecond)
	cancel()
	time.Sleep(100 * time.Millisecond)

	// Validate that the server stopped cleanly.
	connections := server.GetAllConnections()
	assert.Empty(t, connections, "Expected all PeerConnections to be cleaned up")
}

func TestWebRTCServer_DataChannel(t *testing.T) {
	sessionID := "test-session"
	message := "test-message"

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	sdpRequests := make(chan relayerwebrtc.SDPRequest, 1)

	// Test response received from data channel.
	onMessage := func(sessionID, receivedMessage string) {
		assert.Equal(t, message, receivedMessage, "Message mismatch")
	}

	server, err := relayerwebrtc.New(logger, "stun:stun.l.google.com:19302", nil, sdpRequests)
	assert.NoError(t, err, "Failed to create WebRTC server")
	server.OnMessage = onMessage

	peerConnection1, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	assert.NoError(t, err, "Failed to create PeerConnection 1")

	// Set up a DataChannel
	dataChannel, err := peerConnection1.CreateDataChannel("test-data-channel", nil)
	assert.NoError(t, err, "Failed to create DataChannel")

	// Send message on open channel.
	dataChannel.OnOpen(func() {
		err := dataChannel.SendText(message)
		assert.NoError(t, err)
	})

	// Create and send SDP offer
	offer, err := peerConnection1.CreateOffer(nil)
	assert.NoError(t, err, "Failed to create SDP offer")
	err = peerConnection1.SetLocalDescription(offer)
	assert.NoError(t, err, "Failed to set local description")

	responseChan := make(chan *webrtc.SessionDescription)
	sdpRequests <- relayerwebrtc.SDPRequest{
		SessionID: sessionID,
		Offer:     *peerConnection1.LocalDescription(),
		Response:  responseChan,
	}

	// Start server.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := server.Run(ctx)
		assert.NoError(t, err, "WebRTC server exited with error")
	}()

	// Wait for SDP answer.
	answer := <-responseChan
	assert.NotNil(t, answer, "Expected SDP answer")
	err = peerConnection1.SetRemoteDescription(*answer)
	assert.NoError(t, err, "Failed to set remote description")
}

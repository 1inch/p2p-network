package webrtc_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/pion/webrtc/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pb "github.com/1inch/p2p-network/proto"
	relayerwebrtc "github.com/1inch/p2p-network/relayer/webrtc"
)

func TestWebRTCServer_HandleSDP(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	sdpRequests := make(chan relayerwebrtc.SDPRequest, 1)
	iceCandidates := make(chan relayerwebrtc.ICECandidate)

	mockGRPCClient := &mockGRPCClient{}
	mockGRPCClient.On("Close").Return(nil)

	server, err := relayerwebrtc.New(logger, "stun:stun.l.google.com:19302", mockGRPCClient, sdpRequests, iceCandidates)
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
	iceCandidates := make(chan relayerwebrtc.ICECandidate)

	mockGRPCClient := &mockGRPCClient{}
	mockGRPCClient.On("Close").Return(nil)

	server, err := relayerwebrtc.New(logger, "stun:stun.l.google.com:19302", mockGRPCClient, sdpRequests, iceCandidates)
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

	mockGRPCClient.AssertCalled(t, "Close")
}

func TestWebRTCServer_Run_Shutdown(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	sdpRequests := make(chan relayerwebrtc.SDPRequest, 1)
	iceCandidates := make(chan relayerwebrtc.ICECandidate)

	mockGRPCClient := &mockGRPCClient{}
	mockGRPCClient.On("Close").Return(nil)

	server, err := relayerwebrtc.New(logger, "stun:stun.l.google.com:19302", mockGRPCClient, sdpRequests, iceCandidates)
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

	mockGRPCClient.AssertCalled(t, "Close")
}

func TestWebRTCServer_DataChannel(t *testing.T) {
	sessionID := "test-session"
	reqID := "test-req"
	message := "test-message"

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	sdpRequests := make(chan relayerwebrtc.SDPRequest, 1)
	iceCandidates := make(chan relayerwebrtc.ICECandidate)

	mockGRPCClient := &mockGRPCClient{}
	mockGRPCClient.On("Execute", mock.Anything, mock.Anything).Return(&pb.ResolverResponse{
		Id:      "test-id",
		Payload: []byte(message),
		Status:  pb.ResolverResponseStatus_RESOLVER_OK,
	}, nil)

	server, err := relayerwebrtc.New(logger, "stun:stun.l.google.com:19302", mockGRPCClient, sdpRequests, iceCandidates)
	assert.NoError(t, err, "Failed to create WebRTC server")

	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	assert.NoError(t, err, "Failed to create PeerConnection 1")

	// Set up a DataChannel
	dataChannel, err := peerConnection.CreateDataChannel("test-data-channel", nil)
	assert.NoError(t, err, "Failed to create DataChannel")

	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			iceCandidates <- relayerwebrtc.ICECandidate{
				SessionID: sessionID,
				Candidate: *candidate,
			}
		}
	})

	req := &pb.ResolverRequest{
		Id:      reqID,
		Payload: []byte(message),
	}
	reqBytes, err := json.Marshal(req)
	assert.NoError(t, err, "Failed to marshal ResolverRequest")

	respChan := make(chan string, 1)

	// Send message on open channel.
	dataChannel.OnOpen(func() {
		err := dataChannel.SendText(string(reqBytes))
		assert.NoError(t, err)
	})

	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		respChan <- string(msg.Data)
	})

	// Create and send SDP offer
	offer, err := peerConnection.CreateOffer(nil)
	assert.NoError(t, err, "Failed to create SDP offer")
	err = peerConnection.SetLocalDescription(offer)
	assert.NoError(t, err, "Failed to set local description")

	responseChan := make(chan *webrtc.SessionDescription)

	sdpRequests <- relayerwebrtc.SDPRequest{
		SessionID: sessionID,
		Offer:     *peerConnection.LocalDescription(),
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
	err = peerConnection.SetRemoteDescription(*answer)
	assert.NoError(t, err, "Failed to set remote description")

	select {
	case response := <-respChan:
		var resp pb.ResolverResponse
		err := json.Unmarshal([]byte(response), &resp)
		assert.NoError(t, err, "Failed to unmarshal response")
		assert.Equal(t, []byte(message), resp.Payload, "Response payload does not match")
	}

	mockGRPCClient.AssertCalled(t, "Execute", mock.Anything, mock.MatchedBy(func(req *pb.ResolverRequest) bool {
		return req.Id == reqID && string(req.Payload) == message
	}))
}

type mockGRPCClient struct {
	mock.Mock
}

func (m *mockGRPCClient) Execute(ctx context.Context, req *pb.ResolverRequest) (*pb.ResolverResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pb.ResolverResponse), args.Error(1)
}

func (m *mockGRPCClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

package webrtc_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/pion/webrtc/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"

	mocks "github.com/1inch/p2p-network/internal/mock"
	pb "github.com/1inch/p2p-network/proto"
	"github.com/1inch/p2p-network/relayer/grpc"
	relayerwebrtc "github.com/1inch/p2p-network/relayer/webrtc"
)

func TestWebRTCServer_HandleSDP(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	sdpRequests := make(chan relayerwebrtc.SDPRequest, 1)
	iceCandidates := make(chan relayerwebrtc.ICECandidate)

	ctrl := gomock.NewController(t)
	mockGRPCClient := mocks.NewMockGRPCClient(ctrl)
	mockGRPCClient.EXPECT().Close().AnyTimes()

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

	ctrl := gomock.NewController(t)
	mockGRPCClient := mocks.NewMockGRPCClient(ctrl)
	mockGRPCClient.EXPECT().Close().Return(nil)

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
}

func TestWebRTCServer_Run_Shutdown(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	sdpRequests := make(chan relayerwebrtc.SDPRequest, 1)
	iceCandidates := make(chan relayerwebrtc.ICECandidate)

	ctrl := gomock.NewController(t)
	mockGRPCClient := mocks.NewMockGRPCClient(ctrl)
	mockGRPCClient.EXPECT().Close().Return(nil)

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
}

func TestWebRTCServer_DataChannel(t *testing.T) {
	testCases := []struct {
		description string
		setupMock   func(mockGRPCClient *mocks.MockGRPCClient)
		expected    struct {
			errorCode pb.ErrorCode
			errorMsg  string
		}
		requestPayload string
	}{
		{
			description: "Successful gRPC execution",
			setupMock: func(mockGRPCClient *mocks.MockGRPCClient) {
				mockGRPCClient.EXPECT().Execute(gomock.Any(), []byte("public-key"), gomock.Cond(func(req *pb.ResolverRequest) bool {
					return req.Id == "test-req" && string(req.Payload) == "test-message"
				})).Return(&pb.ResolverResponse{
					Id:      "test-id",
					Payload: []byte("test-response"),
				}, nil)
				mockGRPCClient.EXPECT().Close().AnyTimes()
			},
			expected: struct {
				errorCode pb.ErrorCode
				errorMsg  string
			}{
				errorCode: 0,
				errorMsg:  "",
			},
			requestPayload: "test-message",
		},
		{
			description: "Resolver lookup failure",
			setupMock: func(mockGRPCClient *mocks.MockGRPCClient) {
				mockGRPCClient.EXPECT().Execute(gomock.Any(), []byte("public-key"), gomock.Cond(func(req *pb.ResolverRequest) bool {
					return req.Id == "test-req" && string(req.Payload) == "resolver-error"
				})).Return(nil, grpc.ErrResolverLookupFailed)
				mockGRPCClient.EXPECT().Close().AnyTimes()
			},
			expected: struct {
				errorCode pb.ErrorCode
				errorMsg  string
			}{
				errorCode: pb.ErrorCode_ERR_RESOLVER_LOOKUP_FAILED,
				errorMsg:  "resolver lookup failed",
			},
			requestPayload: "resolver-error",
		},
		{
			description: "gRPC execution failure",
			setupMock: func(mockGRPCClient *mocks.MockGRPCClient) {
				mockGRPCClient.EXPECT().Execute(gomock.Any(), []byte("public-key"), gomock.Cond(func(req *pb.ResolverRequest) bool {
					return req.Id == "test-req" && string(req.Payload) == "grpc-error"
				})).Return(nil, grpc.ErrGRPCExecutionFailed)
				mockGRPCClient.EXPECT().Close().AnyTimes()
			},
			expected: struct {
				errorCode pb.ErrorCode
				errorMsg  string
			}{
				errorCode: pb.ErrorCode_ERR_GRPC_EXECUTION_FAILED,
				errorMsg:  "grpc execution failed",
			},
			requestPayload: "grpc-error",
		},
		{
			description: "Incorrect message format",
			setupMock: func(mockGRPCClient *mocks.MockGRPCClient) {
				mockGRPCClient.EXPECT().Execute(gomock.Any(), []byte("public-key"), gomock.Any()).Times(0)
				mockGRPCClient.EXPECT().Close().AnyTimes()
			},
			expected: struct {
				errorCode pb.ErrorCode
				errorMsg  string
			}{
				errorCode: pb.ErrorCode_ERR_INVALID_MESSAGE_FORMAT,
				errorMsg:  "failed to unmarshal protobuf message",
			},
			requestPayload: "invalid-protobuf",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			publicKey := []byte("public-key")
			sessionID := "test-session"
			reqID := "test-req"
			req := &pb.IncomingMessage{
				Request: &pb.ResolverRequest{
					Id:        reqID,
					Payload:   []byte(tc.requestPayload),
					PublicKey: publicKey,
				},
				PublicKeys: [][]byte{
					publicKey,
				},
			}

			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
			sdpRequests := make(chan relayerwebrtc.SDPRequest, 1)
			iceCandidates := make(chan relayerwebrtc.ICECandidate)

			ctrl := gomock.NewController(t)
			mockGRPCClient := mocks.NewMockGRPCClient(ctrl)
			tc.setupMock(mockGRPCClient)

			server, err := relayerwebrtc.New(logger, "stun:stun.l.google.com:19302", mockGRPCClient, sdpRequests, iceCandidates)
			assert.NoError(t, err, "Failed to create WebRTC server")

			peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
			assert.NoError(t, err, "Failed to create PeerConnection")

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

			reqBytes, err := proto.Marshal(req)
			assert.NoError(t, err, "Failed to marshal ResolverRequest")
			if tc.description == "Incorrect message format" {
				reqBytes = []byte(tc.requestPayload)
			}

			respChan := make(chan string, 1)

			dataChannel.OnOpen(func() {
				err := dataChannel.SendText(string(reqBytes))
				assert.NoError(t, err)
			})

			dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
				respChan <- string(msg.Data)
			})

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

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func() {
				err := server.Run(ctx)
				assert.NoError(t, err, "WebRTC server exited with error")
			}()

			answer := <-responseChan
			assert.NotNil(t, answer, "Expected SDP answer")
			err = peerConnection.SetRemoteDescription(*answer)
			assert.NoError(t, err, "Failed to set remote description")

			response := <-respChan
			var resp pb.OutgoingMessage
			err = proto.Unmarshal([]byte(response), &resp)
			assert.NoError(t, err, "Failed to unmarshal response")

			if tc.expected.errorMsg == "" {
				if result, ok := resp.Result.(*pb.OutgoingMessage_Response); ok {
					assert.Equal(t, "test-response", string(result.Response.Payload), "Response payload does not match")
				} else {
					t.Fatalf("Expected a Response result, but got an error: %v", resp.Result)
				}
			} else {
				if result, ok := resp.Result.(*pb.OutgoingMessage_Error); ok {
					assert.Equal(t, tc.expected.errorCode, result.Error.Code, "Error code does not match")
					assert.Contains(t, result.Error.Message, tc.expected.errorMsg, "Error message does not match")
				} else {
					t.Fatalf("Expected an Error result, but got a response: %v", resp.Result)
				}
			}
		})
	}
}

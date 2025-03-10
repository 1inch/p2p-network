package webrtc_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/pion/webrtc/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"

	mocks "github.com/1inch/p2p-network/internal/mock"
	pbrelayer "github.com/1inch/p2p-network/proto/relayer"
	pbresolver "github.com/1inch/p2p-network/proto/resolver"
	"github.com/1inch/p2p-network/relayer/grpc"
	relayerwebrtc "github.com/1inch/p2p-network/relayer/webrtc"
)

var iceServers = []webrtc.ICEServer{{URLs: []string{"stun:stun.l.google.com:19302"}}}

func TestWebRTCServer_HandleSDP(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	sdpRequests := make(chan relayerwebrtc.SDPRequest, 1)
	iceCandidates := make(chan relayerwebrtc.ICECandidate)

	ctrl := gomock.NewController(t)
	mockGRPCClient := mocks.NewMockGRPCClient(ctrl)
	mockGRPCClient.EXPECT().Close().AnyTimes()

	server, err := relayerwebrtc.New(logger, iceServers, mockGRPCClient, sdpRequests, iceCandidates)
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

	server, err := relayerwebrtc.New(logger, iceServers, mockGRPCClient, sdpRequests, iceCandidates)
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

	server, err := relayerwebrtc.New(logger, iceServers, mockGRPCClient, sdpRequests, iceCandidates)
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
	reqID := "test-req"
	testCases := []struct {
		description         string
		setupMock           func(mockGRPCClient *mocks.MockGRPCClient)
		outgoingExpectedErr *struct {
			errorCode pbrelayer.ErrorCode
			errorMsg  string
		}
		expectedPickedPubKey string
		resolverExpectedResp *pbresolver.ResolverResponse
		webrtcOptions        []relayerwebrtc.Option
		// public key generate by this format public-key-<number of key>, example: "public-key-1"
		countPublicKeys uint16
		requestPayload  string
	}{
		{
			description:   "Successful gRPC execution",
			webrtcOptions: []relayerwebrtc.Option{},
			setupMock: func(mockGRPCClient *mocks.MockGRPCClient) {
				mockGRPCClient.EXPECT().Execute(gomock.Any(), []byte("public-key-1"), gomock.Cond(func(req *pbresolver.ResolverRequest) bool {
					return req.Id == "test-req" && string(req.Payload) == "test-message"
				})).Return(&pbresolver.ResolverResponse{
					Id: reqID,
					Result: &pbresolver.ResolverResponse_Payload{
						Payload: []byte("test-response"),
					},
				}, nil)
				mockGRPCClient.EXPECT().Close().AnyTimes()
			},
			outgoingExpectedErr:  nil,
			expectedPickedPubKey: "public-key-1",
			resolverExpectedResp: &pbresolver.ResolverResponse{
				Id: reqID,
				Result: &pbresolver.ResolverResponse_Payload{
					Payload: []byte("test-response"),
				},
			},
			countPublicKeys: 1,
			requestPayload:  "test-message",
		},
		{
			description:   "Error from call method 'execute'",
			webrtcOptions: []relayerwebrtc.Option{},
			setupMock: func(mockGRPCClient *mocks.MockGRPCClient) {
				mockGRPCClient.EXPECT().Execute(gomock.Any(), []byte("public-key-1"), gomock.Cond(func(req *pbresolver.ResolverRequest) bool {
					return req.Id == reqID && string(req.Payload) == "resolver-error"
				})).Return(nil, errors.New("some error from grpc client"))
				mockGRPCClient.EXPECT().Close().AnyTimes()
			},
			outgoingExpectedErr: &struct {
				errorCode pbrelayer.ErrorCode
				errorMsg  string
			}{
				errorCode: pbrelayer.ErrorCode_ERR_GRPC_EXECUTION_FAILED,
				errorMsg:  "failed call execute: some error from grpc client",
			},
			expectedPickedPubKey: "public-key-1",
			resolverExpectedResp: nil,
			countPublicKeys:      1,
			requestPayload:       "resolver-error",
		},
		{
			description:   "Resolver return error 'incorrect message format'",
			webrtcOptions: []relayerwebrtc.Option{},
			setupMock: func(mockGRPCClient *mocks.MockGRPCClient) {
				mockGRPCClient.EXPECT().Execute(gomock.Any(), []byte("public-key-1"), gomock.Any()).
					Times(1).
					Return(&pbresolver.ResolverResponse{
						Id: reqID,
						Result: &pbresolver.ResolverResponse_Error{
							Error: &pbresolver.Error{
								Code:    pbresolver.ErrorCode_ERR_INVALID_MESSAGE_FORMAT,
								Message: "incorrect message",
							},
						},
					}, nil)
				mockGRPCClient.EXPECT().Close().AnyTimes()
			},
			outgoingExpectedErr:  nil,
			expectedPickedPubKey: "public-key-1",
			resolverExpectedResp: &pbresolver.ResolverResponse{
				Id: reqID,
				Result: &pbresolver.ResolverResponse_Error{
					Error: &pbresolver.Error{
						Code:    pbresolver.ErrorCode_ERR_INVALID_MESSAGE_FORMAT,
						Message: "incorrect message",
					},
				},
			},
			countPublicKeys: 1,
			requestPayload:  "invalid-protobuf",
		},
		{
			description: "Retry get request after error",
			webrtcOptions: []relayerwebrtc.Option{
				relayerwebrtc.WithRetry(
					relayerwebrtc.Retry{
						Count: 2,
					},
				),
			},
			setupMock: func(mockGRPCClient *mocks.MockGRPCClient) {
				gomock.InOrder(
					mockGRPCClient.EXPECT().
						Execute(gomock.Any(), []byte("public-key-1"), gomock.Any()).
						Times(1).
						Return(nil, grpc.ErrGRPCExecutionFailed),

					mockGRPCClient.EXPECT().
						Execute(gomock.Any(), []byte("public-key-1"), gomock.Any()).
						Times(1).
						Return(&pbresolver.ResolverResponse{
							Id: reqID,
							Result: &pbresolver.ResolverResponse_Payload{
								Payload: []byte("test-response"),
							},
						}, nil),
				)
				mockGRPCClient.EXPECT().Close().AnyTimes()
			},
			outgoingExpectedErr:  nil,
			expectedPickedPubKey: "public-key-1",
			resolverExpectedResp: &pbresolver.ResolverResponse{
				Id: reqID,
				Result: &pbresolver.ResolverResponse_Payload{
					Payload: []byte("test-response"),
				},
			},
			countPublicKeys: 1,
			requestPayload:  "retry-requests",
		},
		{
			description: "Receive many response one is correct",
			setupMock: func(mockGRPCClient *mocks.MockGRPCClient) {
				mockGRPCClient.EXPECT().
					Execute(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, resolverPublicKey []byte, req *pbresolver.ResolverRequest) (*pbresolver.ResolverResponse, error) {
						if string(resolverPublicKey) == "public-key-3" {
							return &pbresolver.ResolverResponse{
								Id: req.Id,
								Result: &pbresolver.ResolverResponse_Payload{
									Payload: []byte("test-response"),
								},
							}, nil
						}
						return &pbresolver.ResolverResponse{
							Id: req.Id,
						}, grpc.ErrGRPCExecutionFailed
					}).
					// because countPublicKeys is 3 where 1 - will return success, 2 - return errors, retry count 3, so:
					// return errors * retry count + return success = 2 * 3 + 1 = 7 times
					MaxTimes(7)
				mockGRPCClient.EXPECT().Close().AnyTimes()
			},
			webrtcOptions: []relayerwebrtc.Option{
				relayerwebrtc.WithRetry(relayerwebrtc.Retry{
					Count:    3,
					Interval: time.Second,
				})},
			outgoingExpectedErr:  nil,
			expectedPickedPubKey: "public-key-3",
			resolverExpectedResp: &pbresolver.ResolverResponse{
				Id: reqID,
				Result: &pbresolver.ResolverResponse_Payload{
					Payload: []byte("test-response"),
				},
			},
			countPublicKeys: 3,
			requestPayload:  "many-requests",
		},
		{
			description: "Receive many requests all responses error and one final response",
			setupMock: func(mockGRPCClient *mocks.MockGRPCClient) {
				mockGRPCClient.EXPECT().
					Execute(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, resolverPublicKey []byte, req *pbresolver.ResolverRequest) (*pbresolver.ResolverResponse, error) {
						if string(resolverPublicKey) == "public-key-2" {
							return &pbresolver.ResolverResponse{
								Id: req.Id,
								Result: &pbresolver.ResolverResponse_Error{
									Error: &pbresolver.Error{ // if this error returned webrtc cant try another attempt send request
										Code:    pbresolver.ErrorCode_ERR_INVALID_MESSAGE_FORMAT,
										Message: "invalid message",
									},
								},
							}, nil
						}
						return nil, grpc.ErrGRPCExecutionFailed
					}).
					// because countPublicKeys is 3 where 1 - will return success, 2 - return errors, retry count 3, so:
					// return errors * retry count + return success = 2 * 3 + 1 = 7 times
					MaxTimes(7)
				mockGRPCClient.EXPECT().Close().AnyTimes()
			},
			webrtcOptions: []relayerwebrtc.Option{
				relayerwebrtc.WithRetry(relayerwebrtc.Retry{
					Count:    3,
					Interval: time.Second,
				})},
			outgoingExpectedErr: nil,
			resolverExpectedResp: &pbresolver.ResolverResponse{
				Id: reqID,
				Result: &pbresolver.ResolverResponse_Error{
					Error: &pbresolver.Error{ // if this error returned webrtc cant try another attempt send request
						Code:    pbresolver.ErrorCode_ERR_INVALID_MESSAGE_FORMAT,
						Message: "invalid message",
					},
				},
			},
			expectedPickedPubKey: "public-key-2",
			countPublicKeys:      3,
			requestPayload:       "many-requests-return-errors",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			publicKeys := make([][]byte, tc.countPublicKeys)
			for indexPublicKey := range publicKeys {
				publicKeys[indexPublicKey] = []byte(fmt.Sprintf("public-key-%d", indexPublicKey+1))
			}
			req := &pbrelayer.IncomingMessage{
				Request: &pbresolver.ResolverRequest{
					Id:        reqID,
					Payload:   []byte(tc.requestPayload),
					PublicKey: []byte("public-key"),
				},
				PublicKeys: publicKeys,
			}

			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))
			sdpRequests := make(chan relayerwebrtc.SDPRequest, 1)
			iceCandidates := make(chan relayerwebrtc.ICECandidate)

			ctrl := gomock.NewController(t)
			mockGRPCClient := mocks.NewMockGRPCClient(ctrl)
			tc.setupMock(mockGRPCClient)

			server, err := relayerwebrtc.New(logger, iceServers, mockGRPCClient, sdpRequests, iceCandidates, tc.webrtcOptions...)

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
			var outgoingResponseMessage pbrelayer.OutgoingMessage
			err = proto.Unmarshal([]byte(response), &outgoingResponseMessage)
			assert.NoError(t, err, "Failed to unmarshal response")

			if tc.outgoingExpectedErr != nil {
				respErr := outgoingResponseMessage.GetError()
				assert.NotNil(t, respErr, "Expected error in 'OutgoinMessage")

				assert.Equal(t, tc.outgoingExpectedErr.errorCode, respErr.Code, "error code in OutgoingMessage doesnt match")
				assert.Equal(t, tc.outgoingExpectedErr.errorMsg, respErr.Message, "error message in OutgoingMessage doesnt match")
			} else {
				assert.Equal(t, []byte(tc.expectedPickedPubKey), outgoingResponseMessage.PublicKey)

				resolverResponse := outgoingResponseMessage.GetResponse()

				assert.Equal(t, tc.resolverExpectedResp, resolverResponse)

			}
		})
	}
}

// func TestRetry(t *testing.T) {
// 	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
// 		Level: slog.LevelDebug,
// 	}))
// 	sdpRequests := make(chan relayerwebrtc.SDPRequest, 1)
// 	iceCandidates := make(chan relayerwebrtc.ICECandidate)

// 	conn, _ := grpc2.NewClient("127.0.0.1:8001",
// 		grpc2.WithTransportCredentials(insecure.NewCredentials()),
// 	)

// 	conns := make(map[string]*grpc2.ClientConn)

// 	pubkey := ""
// 	id := ""

// 	conns[pubkey] = conn

// 	grpcClient := grpc.New2(
// 		logger.WithGroup("grpc-server"),
// 		conns)

// 	server, _ := relayerwebrtc.New(logger, iceServers, grpcClient, sdpRequests, iceCandidates)

// 	doneChan := make(chan bool)
// 	respChan := make(chan *pbrelayer.OutgoingMessage)
// 	server.RetryGetResponseFromResolver([]byte(pubkey), &pbresolver.ResolverRequest{
// 		Id:        id,
// 		Payload:   []byte{},
// 		PublicKey: []byte{},
// 	}, doneChan, respChan)

// 	<-doneChan
// 	resp := <-respChan

// 	logger.Info(string(resp.PublicKey))
// }

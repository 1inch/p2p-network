package test

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"testing"

	// "github.com/1inch/p2p-network/relayer"

	pb "github.com/1inch/p2p-network/proto"
	relayergrpc "github.com/1inch/p2p-network/relayer/grpc"
	relayerwebrtc "github.com/1inch/p2p-network/relayer/webrtc"
	"github.com/1inch/p2p-network/resolver"
	"github.com/1inch/p2p-network/resolver/types"
	"github.com/pion/webrtc/v4"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	// "google.golang.org/grpc"
)

const (
	httpEndpointToRelayer  = "127.0.0.1:8080"
	grpcEndpointToResolver = "127.0.0.1:8001"
	ICEServer              = "stun:stun1.l.google.com:19302"
)

type positiveTestCase struct {
	Name                              string
	SessionId                         string
	JsonRequest                       *types.JsonRequest
	FuncCheckActualJsonResponseResult func(result interface{})
	ConfigForResolver                 *resolver.Config
}

func TestPositiveCases(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// maybe need use some pseudo-random algorithm for generate session_id, requestId,
	testCases := []positiveTestCase{
		{
			Name:      "ResolverUsedDefaultHandler",
			SessionId: "test-session-id-1",
			JsonRequest: &types.JsonRequest{
				Id:     "request-id-1",
				Method: "GetWalletBalance",
				Params: []string{"0x0ADfCCa4B2a1132F82488546AcA086D7E24EA324", "latest"},
			},
			FuncCheckActualJsonResponseResult: func(result interface{}) {
				assert.Equal(t, 555., result)
			},
			ConfigForResolver: cfgResolverWithDefaultApi(),
		},
		{
			Name:      "ResolverUsedInfuraHandler",
			SessionId: "test-session-id-2",
			JsonRequest: &types.JsonRequest{
				Id:     "request-id-2",
				Method: "GetWalletBalance",
				Params: []string{"0x0ADfCCa4B2a1132F82488546AcA086D7E24EA324", "latest"},
			},
			FuncCheckActualJsonResponseResult: func(result interface{}) {
				assert.NotEmpty(t, result)
			},
			ConfigForResolver: cfgResolverWithInfuraApi(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			positiveTestWorkFlow(t, logger, &testCase)
		})
	}
}

func TestUnhandledMessageInResolver(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	expectedCode := codes.InvalidArgument
	expecteErrorMessage := "unrecognized method"
	jsonReq := &types.JsonRequest{
		Id:     "request-id-2",
		Method: "blockNumber",
	}
	respBytes := testWorkFlowAndReturnResponseChan(t, logger, cfgResolverWithDefaultApi(),
		"test-session-id-3",
		jsonReq,
	)

	var errorDetails error
	err := json.Unmarshal(respBytes, &errorDetails)
	assert.NoError(t, err, "Failed to unmarshal response")

	statusError := status.Convert(errorDetails)
	responseResolver, ok := statusError.Details()[0].(*pb.ResolverResponse)

	assert.True(t, ok, "expect that first elem in status error detais is responseResolver")

	assert.Equal(t, jsonReq.Id, responseResolver.Id)
	assert.Equal(t, expectedCode, statusError.Code())
	assert.Equal(t, expecteErrorMessage, statusError.Message())
}

func testWorkFlowAndReturnResponseChan(t *testing.T, logger *slog.Logger, cfg *resolver.Config,
	sessionId string, jsonReq *types.JsonRequest) []byte {

	sdpRequests, iceCandidates, relayerServer, err := setupRelayer(logger)
	assert.NoError(t, err, "Failed to create WebRTC server")

	resolverGrpcServer, err := setupResolver(cfg)
	assert.NoError(t, err, "Failed to create resolver grpc server")
	defer resolverGrpcServer.GracefulStop()

	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	assert.NoError(t, err, "Failed to create PeerConnection 1")

	// Set up a DataChannel
	dataChannel, err := peerConnection.CreateDataChannel("test-data-channel", nil)
	assert.NoError(t, err, "Failed to create DataChannel")

	// change this to call relayer "candidate" endpoint when fix session
	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			iceCandidates <- relayerwebrtc.ICECandidate{
				SessionID: sessionId,
				Candidate: *candidate,
			}
		}
	})

	payload, err := json.Marshal(jsonReq)
	assert.NoError(t, err, "Failed to marshal JsonRequest")

	req := &pb.ResolverRequest{
		Id:      jsonReq.Id,
		Payload: payload,
	}
	reqBytes, err := json.Marshal(req)
	assert.NoError(t, err, "Failed to marshal ResolverRequest")

	respChan := make(chan []byte, 1)

	dataChannel.OnOpen(func() {
		err := dataChannel.Send(reqBytes)
		assert.NoError(t, err)
	})

	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		respChan <- msg.Data
	})

	// Create and send SDP offer
	offer, err := peerConnection.CreateOffer(nil)
	assert.NoError(t, err, "Failed to create SDP offer")
	err = peerConnection.SetLocalDescription(offer)
	assert.NoError(t, err, "Failed to set local description")

	responseChan := make(chan *webrtc.SessionDescription)

	// change this to call relayer "sdp" endpoint when fix session
	sdpRequests <- relayerwebrtc.SDPRequest{
		SessionID: sessionId,
		Offer:     *peerConnection.LocalDescription(),
		Response:  responseChan,
	}

	// Start relayer server.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := relayerServer.Run(ctx)
		assert.NoError(t, err, "WebRTC server exited with error")
	}()

	// Wait for SDP answer.
	answer := <-responseChan
	assert.NotNil(t, answer, "Expected SDP answer")
	err = peerConnection.SetRemoteDescription(*answer)
	assert.NoError(t, err, "Failed to set remote description")

	return <-respChan
}

func positiveTestWorkFlow(t *testing.T, logger *slog.Logger, testCase *positiveTestCase) {
	respBytes := testWorkFlowAndReturnResponseChan(t, logger, testCase.ConfigForResolver, testCase.SessionId, testCase.JsonRequest)
	var resp pb.ResolverResponse
	err := json.Unmarshal(respBytes, &resp)
	assert.NoError(t, err, "Failed to unmarshal response")

	var jsonResp types.JsonResponse
	err = json.Unmarshal(resp.Payload, &jsonResp)
	assert.NoError(t, err, "Failed to unmarshal response")

	testCase.FuncCheckActualJsonResponseResult(jsonResp.Result)
}

func setupRelayer(logger *slog.Logger) (chan relayerwebrtc.SDPRequest, chan relayerwebrtc.ICECandidate, *relayerwebrtc.Server, error) {
	sdpRequests := make(chan relayerwebrtc.SDPRequest, 1)
	iceCandidates := make(chan relayerwebrtc.ICECandidate)
	grpcClient, err := relayergrpc.New(grpcEndpointToResolver)
	if err != nil {
		return nil, nil, nil, err
	}

	server, err := relayerwebrtc.New(logger, ICEServer, grpcClient, sdpRequests, iceCandidates)
	if err != nil {
		return nil, nil, nil, err
	}

	return sdpRequests, iceCandidates, server, nil
}

func setupResolver(cfg *resolver.Config) (*grpc.Server, error) {
	return resolver.Run(cfg)
}

func cfgResolverWithDefaultApi() *resolver.Config {
	cfgResolver := cfgResolverWithoutApis()

	cfgResolver.Apis = resolver.ApiConfigs{
		Default: resolver.DefaultApiConfig{
			Enabled: true,
		},
	}

	return cfgResolver
}

func cfgResolverWithInfuraApi() *resolver.Config {
	cfgResolver := cfgResolverWithoutApis()

	cfgResolver.Apis = resolver.ApiConfigs{
		Infura: resolver.InfuraApiConfig{
			Enabled: true,
			Key:     "a8401733346d412389d762b5a63b0bcf",
		},
	}

	return cfgResolver
}

func cfgResolverWithoutApis() *resolver.Config {
	return &resolver.Config{
		Port:     8001,
		LogLevel: slog.LevelInfo,
	}
}

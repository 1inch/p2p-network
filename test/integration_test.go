package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	// "github.com/1inch/p2p-network/relayer"

	pb "github.com/1inch/p2p-network/proto"
	"github.com/1inch/p2p-network/relayer"
	"github.com/1inch/p2p-network/resolver"
	"github.com/1inch/p2p-network/resolver/types"
	"github.com/pion/webrtc/v4"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	// "google.golang.org/grpc"
)

const (
	formatForUrl           = "http://%s/%s"
	httpEndpointToRelayer  = "127.0.0.1:8080"
	grpcEndpointToResolver = "127.0.0.1:8001"
	ICEServer              = "stun:stun1.l.google.com:19302"
)

type TestCase struct {
	Name                              string
	SessionId                         string
	JsonRequest                       *types.JsonRequest
	ExpectedResolverApiName           string
	ExpectedResolverResponseStatus    pb.ResolverResponseStatus
	ExpectedJsonResponseError         string
	FuncCheckActualJsonResponseResult func(result interface{})
	ConfigForResolver                 *resolver.Config
}

func TestRelayerAndResolverIntegration(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// maybe need use some pseudo-random algorithm for generate session_id, requestId,
	testCases := []TestCase{
		{
			Name:      "ResolverUsedDefaultHandler",
			SessionId: "test-session-id-1",
			JsonRequest: &types.JsonRequest{
				Id:     "request-id-1",
				Method: "GetWalletBalance",
				Params: []string{"0x0ADfCCa4B2a1132F82488546AcA086D7E24EA324", "latest"},
			},
			ExpectedResolverApiName:        "default",
			ExpectedResolverResponseStatus: pb.ResolverResponseStatus_RESOLVER_OK,
			ExpectedJsonResponseError:      "",
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
			ExpectedResolverApiName:        "infura",
			ExpectedResolverResponseStatus: pb.ResolverResponseStatus_RESOLVER_OK,
			ExpectedJsonResponseError:      "",
			FuncCheckActualJsonResponseResult: func(result interface{}) {
				assert.NotEmpty(t, result)
			},
			ConfigForResolver: cfgResolverWithInfuraApi(),
		},
		{
			Name:      "UnhandledMessageInResolver",
			SessionId: "test-session-id-3",
			JsonRequest: &types.JsonRequest{
				Id:     "request-id-2",
				Method: "blockNumber",
			},
			ExpectedResolverApiName:           "default",
			ExpectedResolverResponseStatus:    pb.ResolverResponseStatus_RESOLVER_ERROR, // this status code is correct in this case ?
			ExpectedJsonResponseError:         "Unrecognized method",
			FuncCheckActualJsonResponseResult: func(result interface{}) { assert.Equal(t, float64(0), result) },
			ConfigForResolver:                 cfgResolverWithDefaultApi(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			testWorkFlow(t, logger, &testCase)
		})
	}
}

func testWorkFlow(t *testing.T, logger *slog.Logger, testCase *TestCase) {
	time.Sleep(time.Duration(1000)) // add this sleep because sometimes http server does not have time to stop before start new test
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := setupRelayer(ctx, logger)
	assert.NoError(t, err, "Failed to create WebRTC server")

	resolverGrpcServer, err := setupResolver(testCase.ConfigForResolver)
	assert.NoError(t, err, "Failed to create resolver grpc server")
	defer resolverGrpcServer.GracefulStop()

	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	assert.NoError(t, err, "Failed to create PeerConnection 1")

	// Set up a DataChannel
	dataChannel, err := peerConnection.CreateDataChannel("test-data-channel", nil)
	assert.NoError(t, err, "Failed to create DataChannel")

	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			logger.Info("start send candidate", slog.Any("address", candidate.Address))
			resp, err := sendCandidateRequestToRelayer(testCase.SessionId, candidate, logger)
			defer func() {
				err := resp.Body.Close()
				if err != nil {
					logger.Error("error when try close response body")
				}
			}()
			if err != nil {
				logger.Error("error when try send candidate request to relayer", slog.Any("err", err.Error()))
			}
			if resp.StatusCode != http.StatusAccepted {
				logger.Error("status code from candidate response not equal \"200 OK\"")
			}
		}
	})

	payload, err := json.Marshal(testCase.JsonRequest)
	assert.NoError(t, err, "Failed to marshal JsonRequest")

	req := &pb.ResolverRequest{
		Id:      testCase.JsonRequest.Id,
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

	sdpResponse := struct {
		Answer webrtc.SessionDescription `json:"answer"`
	}{}

	logger.Info("start send sdp")
	httpSdpResp, err := sendSDPRequestToRelayer(testCase.SessionId, peerConnection.LocalDescription(), logger)
	defer func() {
		err := httpSdpResp.Body.Close()
		if err != nil {
			logger.Error("error when try close response body")
		}
	}()
	assert.NoError(t, err, "Failed to send sdp request to relayer")
	assert.Equal(t, http.StatusOK, httpSdpResp.StatusCode, "Expected status ok from sdp endpoint")

	body, err := io.ReadAll(httpSdpResp.Body)
	assert.NoError(t, err, "Failed read body from response")
	err = json.Unmarshal(body, &sdpResponse)
	assert.NoError(t, err, "Failed unmarshal body from response")
	assert.NotNil(t, sdpResponse, "Expected SDP response")

	err = peerConnection.SetRemoteDescription(sdpResponse.Answer)
	assert.NoError(t, err, "Failed to set remote description")

	response := <-respChan
	var resp pb.ResolverResponse
	err = json.Unmarshal(response, &resp)
	assert.NoError(t, err, "Failed to unmarshal response")
	assert.Equal(t, testCase.ExpectedResolverResponseStatus, resp.Status, "Response status not equal")

	var jsonResps types.JsonResponses
	err = json.Unmarshal(resp.Payload, &jsonResps)
	assert.NoError(t, err, "Failed to unmarshal response")

	jsonResp, ok := jsonResps[testCase.ExpectedResolverApiName]
	assert.Truef(t, ok, "Resolver not return response for api name: %s", testCase.ExpectedResolverApiName)
	testCase.FuncCheckActualJsonResponseResult(jsonResp.Result)
	assert.Equal(t, testCase.ExpectedJsonResponseError, jsonResp.Error)
}

func sendCandidateRequestToRelayer(sessionId string, candidate *webrtc.ICECandidate, logger *slog.Logger) (*http.Response, error) {
	req := struct {
		SessionID string              `json:"session_id"`
		Candidate webrtc.ICECandidate `json:"candidate"`
	}{
		SessionID: sessionId,
		Candidate: *candidate,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		logger.Error("error when try marshal ice candidate", slog.Any("err", err.Error()))
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, fmt.Sprintf(formatForUrl, httpEndpointToRelayer, "candidate"), bytes.NewReader(reqBytes))
	if err != nil {
		logger.Error("error when create http request to candidate")
		return nil, err
	}
	return http.DefaultClient.Do(httpReq)
}

func sendSDPRequestToRelayer(sessionId string, sessionDescription *webrtc.SessionDescription, logger *slog.Logger) (*http.Response, error) {
	req := struct {
		SessionID string                    `json:"session_id"`
		Offer     webrtc.SessionDescription `json:"offer"`
	}{
		SessionID: sessionId,
		Offer:     *sessionDescription,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		logger.Error("error when try marshal session description", slog.Any("err", err.Error()))
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, fmt.Sprintf(formatForUrl, httpEndpointToRelayer, "sdp"), bytes.NewReader(reqBytes))
	if err != nil {
		logger.Error("error when create http request to candidate")
		return nil, err
	}
	return http.DefaultClient.Do(httpReq)
}

func setupRelayer(ctx context.Context, logger *slog.Logger) (*relayer.Relayer, error) {
	cfg := &relayer.Config{
		LogLevel:          "INFO",
		HTTPEndpoint:      httpEndpointToRelayer,
		WebRTCICEServer:   ICEServer,
		GRPCServerAddress: grpcEndpointToResolver,
	}
	relayerNode, err := relayer.New(cfg, logger.WithGroup("Relayer"))
	if err != nil {
		return nil, err
	}

	go func() {
		err := relayerNode.Run(ctx)
		if err != nil {
			logger.Error("error when try start relayer node")
		}
	}()

	return relayerNode, nil
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

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

	pb "github.com/1inch/p2p-network/proto"
	"github.com/1inch/p2p-network/relayer"
	"github.com/1inch/p2p-network/resolver"
	"github.com/1inch/p2p-network/resolver/types"
	"github.com/phayes/freeport"
	"github.com/pion/webrtc/v4"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

const (
	formatForUrl                 = "http://%s/%s"
	formatHttpEndpointToRelayer  = "127.0.0.1:%d"
	formatGrpcEndpointToResolver = "127.0.0.1:%d"
	ICEServer                    = "stun:stun1.l.google.com:19302"
	blockchainRPCAddress         = "http://127.0.0.1:8545"
	deploymentPrivateKey         = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	relayerPrivateKey            = "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"
	contractAddress              = "0x5fbdb2315678afecb367f032d93f642f64180aa3"
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
		// {
		// 	Name:      "UnhandledMessageInResolver",
		// 	SessionId: "test-session-id-3",
		// 	JsonRequest: &types.JsonRequest{
		// 		Id:     "request-id-2",
		// 		Method: "blockNumber",
		// 	},
		// 	ExpectedResolverApiName:           "default",
		// 	ExpectedResolverResponseStatus:    pb.ResolverResponseStatus_RESOLVER_ERROR, // this status code is correct in this case ?
		// 	ExpectedJsonResponseError:         "Unrecognized method",
		// 	FuncCheckActualJsonResponseResult: func(result interface{}) { assert.Equal(t, float64(0), result) },
		// 	ConfigForResolver:                 cfgResolverWithDefaultApi(),
		// },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			testWorkFlow(t, logger, &testCase)
		})
	}
}

func testWorkFlow(t *testing.T, logger *slog.Logger, testCase *TestCase) {
	logger.Info("start test case", slog.Any("name", testCase.Name))

	ctx, cancel := context.WithCancel(context.Background())

	// _, err := registry.DeployNodeRegistry(ctx, registry.Config{
	// 	DialURI:    blockchainRPCAddress,
	// 	PrivateKey: deploymentPrivateKey,
	// })
	// assert.NoError(t, err, "Failed to deploy Node Registry contract")

	_, httpRelayerUri, err := setupRelayer(ctx, logger, testCase.ConfigForResolver)
	assert.NoError(t, err, "Failed to create WebRTC server")
	defer cancel()

	resolverGrpcServer, err := setupResolver(testCase.ConfigForResolver)
	assert.NoError(t, err, "Failed to create resolver grpc server")
	defer resolverGrpcServer.GracefulStop() // its block stopped servers and test-case

	// Create new peer connection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	assert.NoError(t, err, "Failed to create PeerConnection 1")

	// Set up a DataChannel
	dataChannel, err := peerConnection.CreateDataChannel("test-data-channel", nil)
	assert.NoError(t, err, "Failed to create DataChannel")

	// Add handler for IceCandidate, in the handler need send this candidate to relayer by http endpoint
	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			logger.Info("start send candidate", slog.Any("address", candidate.Address))
			resp, err := sendCandidateRequestToRelayer(testCase.SessionId, httpRelayerUri, candidate, logger)
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

	// start send sdp request to relayer
	logger.Info("start send sdp")
	httpSdpResp, err := sendSDPRequestToRelayer(testCase.SessionId, httpRelayerUri, peerConnection.LocalDescription(), logger)
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

	// get response from data channel and check the response
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

// sendCandidateRequestToRelayer method for send request candidate to relayer using http endpoint
func sendCandidateRequestToRelayer(sessionId string, relayerUri string, candidate *webrtc.ICECandidate, logger *slog.Logger) (*http.Response, error) {
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

	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, fmt.Sprintf(formatForUrl, relayerUri, "candidate"), bytes.NewReader(reqBytes))
	if err != nil {
		logger.Error("error when create http request to candidate")
		return nil, err
	}
	return http.DefaultClient.Do(httpReq)
}

// sendSDPRequestToRelayer method for send request sdp to relayer using http endpoint
func sendSDPRequestToRelayer(sessionId string, relayerUri string, sessionDescription *webrtc.SessionDescription, logger *slog.Logger) (*http.Response, error) {
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

	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, fmt.Sprintf(formatForUrl, relayerUri, "sdp"), bytes.NewReader(reqBytes))
	if err != nil {
		logger.Error("error when create http request to candidate")
		return nil, err
	}
	return http.DefaultClient.Do(httpReq)
}

func setupRelayer(ctx context.Context, logger *slog.Logger, resolverCfg *resolver.Config) (*relayer.Relayer, string, error) {
	portForHttp, err := freeport.GetFreePort()

	if err != nil {
		return nil, "", err
	}

	httpEndpointToRelayer := fmt.Sprintf(formatHttpEndpointToRelayer, portForHttp)
	grpcEndpointToResolver := fmt.Sprintf(formatGrpcEndpointToResolver, resolverCfg.Port)
	cfg := &relayer.Config{
		LogLevel:             "INFO",
		HTTPEndpoint:         httpEndpointToRelayer,
		WebRTCICEServer:      ICEServer,
		GRPCServerAddress:    grpcEndpointToResolver,
		BlockchainRPCAddress: blockchainRPCAddress,
		PrivateKey:           relayerPrivateKey,
		ContractAddress:      contractAddress,
	}
	relayerNode, err := relayer.New(cfg, logger.WithGroup("Relayer"))
	if err != nil {
		return nil, "", err
	}

	go func() {
		err := relayerNode.Run(ctx)
		if err != nil {
			logger.Error("error when try start relayer node")
		}
	}()

	return relayerNode, httpEndpointToRelayer, nil
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
	port, err := freeport.GetFreePort()

	if err != nil {
		panic("cant get free port for resolver")
	}

	return &resolver.Config{
		Port:     port,
		LogLevel: slog.LevelInfo,
	}
}

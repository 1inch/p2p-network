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

	"github.com/1inch/p2p-network/internal/registry"
	"github.com/1inch/p2p-network/internal/testnetwork"
	pbrelayer "github.com/1inch/p2p-network/proto/relayer"
	pbresolver "github.com/1inch/p2p-network/proto/resolver"
	relayergrpc "github.com/1inch/p2p-network/relayer/grpc"
	relayerwebrtc "github.com/1inch/p2p-network/relayer/webrtc"
	"github.com/1inch/p2p-network/resolver"
	"github.com/1inch/p2p-network/resolver/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pion/webrtc/v4"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

const (
	httpEndpointToRelayer  = "127.0.0.1:8080"
	grpcEndpointToResolver = "127.0.0.1:8001"
	dialURL                = "127.0.0.1:8545"
	privateKey             = "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"
	resolverPrivateKey     = "5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"
	contractAddress        = "0x5fbdb2315678afecb367f032d93f642f64180aa3"
)

var iceServers = []webrtc.ICEServer{{URLs: []string{"stun:stun.l.google.com:19302"}}}

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

// TODO TestUnhandledMessageInResolver uncomment this test when fix error handling on relayer
// func TestUnhandledMessageInResolver(t *testing.T) {
// 	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
// 	expectedCode := codes.InvalidArgument
// 	expecteErrorMessage := "unrecognized method"
// 	jsonReq := &types.JsonRequest{
// 		Id:     "request-id-2",
// 		Method: "blockNumber",
// 	}
// 	respBytes := testWorkFlowAndReturnResponseChan(t, logger, cfgResolverWithDefaultApi(),
// 		"test-session-id-3",
// 		jsonReq,
// 	)

// 	var errorDetails error
// 	err := json.Unmarshal(respBytes, &errorDetails)
// 	assert.NoError(t, err, "Failed to unmarshal response")

// 	statusError := status.Convert(errorDetails)
// 	responseResolver, ok := statusError.Details()[0].(*pb.ResolverResponse)

// 	assert.True(t, ok, "expect that first elem in status error detais is responseResolver")

// 	assert.Equal(t, jsonReq.Id, responseResolver.Id)
// 	assert.Equal(t, expectedCode, statusError.Code())
// 	assert.Equal(t, expecteErrorMessage, statusError.Message())
// }

func testWorkFlowAndReturnResponseChan(t *testing.T, logger *slog.Logger, cfg *resolver.Config,
	sessionId string, jsonReq *types.JsonRequest) []byte {

	sdpRequests, iceCandidates, relayerServer, err := setupRelayer(logger)
	assert.NoError(t, err, "Failed to create WebRTC server")

	resolverGrpcServer, err := setupResolver(cfg, logger)
	assert.NoError(t, err, "Failed to create resolver node")
	err = resolverGrpcServer.Run()
	assert.NoError(t, err, "Failed to start resolver")
	defer func() {
		err = resolverGrpcServer.Stop()
		assert.NoError(t, err, "Failed to stop resolver")
	}()

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

	privKey, err := crypto.HexToECDSA(resolverPrivateKey)
	assert.NoError(t, err, "invalid private key")
	resolverPublicKeyBytes := crypto.CompressPubkey(&privKey.PublicKey)

	req := &pbrelayer.IncomingMessage{
		Request: &pbresolver.ResolverRequest{
			Id:      jsonReq.Id,
			Payload: payload,
		},
		PublicKeys: [][]byte{
			resolverPublicKeyBytes,
		},
	}
	reqBytes, err := proto.Marshal(req)
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
	var resp pbrelayer.OutgoingMessage
	err := proto.Unmarshal(respBytes, &resp)
	assert.NoError(t, err, "Failed to unmarshal response")

	if result, ok := resp.Result.(*pbrelayer.OutgoingMessage_Response); ok {
		var jsonResp types.JsonResponse
		err = json.Unmarshal(result.Response.GetPayload(), &jsonResp)
		assert.NoError(t, err, "Failed to unmarshal response")

		testCase.FuncCheckActualJsonResponseResult(jsonResp.Result)
	}
}

func setupRelayer(logger *slog.Logger) (chan relayerwebrtc.SDPRequest, chan relayerwebrtc.ICECandidate, *relayerwebrtc.Server, error) {
	ctx := context.Background()
	sdpRequests := make(chan relayerwebrtc.SDPRequest, 1)
	iceCandidates := make(chan relayerwebrtc.ICECandidate)
	registry, err := registry.Dial(ctx, &registry.Config{
		DialURI:         fmt.Sprintf("http://%s", dialURL),
		PrivateKey:      privateKey,
		ContractAddress: contractAddress,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	server, err := relayerwebrtc.New(logger, iceServers, relayergrpc.New(logger, registry), sdpRequests, iceCandidates)

	if err != nil {
		return nil, nil, nil, err
	}

	return sdpRequests, iceCandidates, server, nil
}

func setupResolver(cfg *resolver.Config, logger *slog.Logger) (*resolver.Resolver, error) {
	return resolver.New(*cfg, logger)
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
		GrpcEndpoint: grpcEndpointToResolver,
		LogLevel:     slog.LevelInfo,
	}
}

func TestSuccess(t *testing.T) {
	testnetwork.Run(t, 1, 1, func(tn *testnetwork.TestNetwork) {
		ctx := context.Background()
		client := &http.Client{}
		resolverCfg := &resolver.Config{
			RpcUrl:          dialURL,
			GrpcEndpoint:    tn.ResolverNodes[0].Addr(),
			ContractAddress: contractAddress,
			PrivateKey:      tn.ResolverPrivateKeys[1],
		}
		resolverRegistrationService, err := resolver.NewRegistrationResolver(slog.Default(), resolverCfg)
		assert.Nil(t, err, "Failed create registration service for resolver")
		_, err = resolverRegistrationService.Register(ctx)
		assert.Nil(t, err, "Failed registration resolver")
		relayerAddress := tn.RelayerNodes[0].HTTPServer.Addr()
		jsonReq := &types.JsonRequest{
			Id:     "request-id-1",
			Method: "GetWalletBalance",
			Params: []string{"0x0ADfCCa4B2a1132F82488546AcA086D7E24EA324", "latest"},
		}

		peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
		assert.NoError(t, err, "Failed to create PeerConnection 1")

		// Set up a DataChannel
		dataChannel, err := peerConnection.CreateDataChannel("test-data-channel", nil)
		assert.NoError(t, err, "Failed to create DataChannel")

		// change this to call relayer "candidate" endpoint when fix session
		peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
			if candidate != nil {
				payload := map[string]interface{}{
					"session_id": "test-session",
					"candidate":  *candidate,
				}

				p, err := json.Marshal(payload)
				assert.NoError(t, err, "Failed to marshal ICECandidate")

				req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://%s/candidate", relayerAddress), bytes.NewBuffer(p))
				assert.NoError(t, err, "Failed to create POST request")
				req.Header.Set("Content-Type", "application/json")

				b, err := client.Do(req)
				assert.NoError(t, err, "Failed to send POST request")
				err = b.Body.Close()
				assert.NoError(t, err, "Failed to close response body")
			}
		})

		payload, err := json.Marshal(jsonReq)
		assert.NoError(t, err, "Failed to marshal JsonRequest")

		privKey, err := crypto.HexToECDSA(resolverCfg.PrivateKey)
		assert.NoError(t, err, "invalid private key")
		resolverPublicKeyBytes := crypto.CompressPubkey(&privKey.PublicKey)

		msg := &pbrelayer.IncomingMessage{
			Request: &pbresolver.ResolverRequest{
				Id:      jsonReq.Id,
				Payload: payload,
			},
			PublicKeys: [][]byte{
				resolverPublicKeyBytes,
			},
		}
		msgBytes, err := proto.Marshal(msg)
		assert.NoError(t, err, "Failed to marshal ResolverRequest")

		respChan := make(chan []byte, 1)

		dataChannel.OnOpen(func() {
			err := dataChannel.Send(msgBytes)
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

		// Do sdp request
		sdpPayload := map[string]interface{}{
			"session_id": "test-session",
			"offer":      offer,
		}

		sdpReq, err := json.Marshal(sdpPayload)
		assert.NoError(t, err, "Failed to marshal SDP offer")

		req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://%s/sdp", relayerAddress), bytes.NewBuffer(sdpReq))
		assert.NoError(t, err, "Failed to create SDP request")
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send POST request: %v", err)
		}
		t.Cleanup(
			func() {
				err := resp.Body.Close()
				assert.NoError(t, err, "Failed to close response body")
			},
		)

		type SDPResponse struct {
			Answer webrtc.SessionDescription `json:"answer"`
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Unexpected status code: %d. Response body: %s", resp.StatusCode, string(bodyBytes))
		}

		var sdpResp SDPResponse
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&sdpResp); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		t.Logf("Received SDP answer: %v", sdpResp.Answer)
		fmt.Printf("Received SDP answer: %v", sdpResp.Answer)

		err = peerConnection.SetRemoteDescription(sdpResp.Answer)
		assert.NoError(t, err, "Failed to set remote description")

		respBytes := <-respChan
		var outMsg pbrelayer.OutgoingMessage
		err = proto.Unmarshal(respBytes, &outMsg)
		assert.NoError(t, err, "Failed to unmarshal response")

		assert.Nil(t, outMsg.GetError(), "Unexpected error in OutgoingMessage")
		resolverResponse := outMsg.GetResponse()
		assert.Nil(t, resolverResponse.GetError(), "Unexpected error in ResolverResponse")

		var jsonResp types.JsonResponse
		err = json.Unmarshal(resolverResponse.GetPayload(), &jsonResp)
		t.Logf("json response: %s", jsonResp)
		assert.NoError(t, err, "Failed to unmarshal response")

		assert.Equal(t, 555., jsonResp.Result)
	})
}

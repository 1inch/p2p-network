package qa

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

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

	rand.Seed(time.Now().UnixNano())

	testCases := []positiveTestCase{
		{
			Name:      "ResolverUsedOneInchHandler",
			SessionId: generateSessionID(),
			JsonRequest: &types.JsonRequest{
				Id:     generateRequestID(),
				Method: "GetWalletBalance",
				Params: []string{"1", "0x0ADfCCa4B2a1132F82488546AcA086D7E24EA324"},
			},
			FuncCheckActualJsonResponseResult: func(result interface{}) {
				assert.NotEmpty(t, result)
			},
			ConfigForResolver: cfgResolverWithOneInchApi(),
		},
		// GetWalletBalanceForVitalikButerinAddress: This test case checks the wallet balance retrieval for Vitalik Buterin's address using the Infura API. It verifies that the result is a non-empty hex string starting with "0x".
		{
			Name:      "GetWalletBalanceForVitalikButerinAddress",
			SessionId: generateSessionID(),
			JsonRequest: &types.JsonRequest{
				Id:     generateRequestID(),
				Method: "GetWalletBalance",
				Params: []string{"1", "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"}, // Vitalik Buterin's address
			},
			FuncCheckActualJsonResponseResult: func(result interface{}) {
				assert.NotEmpty(t, result)
				// For OneInch, result should be a map of token addresses to balances
				balanceMap, ok := result.(map[string]interface{})
				assert.True(t, ok, "Expected result to be a map")

				// Check that the map is not empty
				assert.NotEmpty(t, balanceMap, "Expected non-empty balance map")

				// Check each key-value pair in the map
				for tokenAddress, balance := range balanceMap {
					// Check that the token address is a valid Ethereum address
					assert.Regexp(t, "^0x[a-fA-F0-9]{40}$", tokenAddress, "Expected token address to be a valid Ethereum address")

					// Check that the balance is a string
					balanceStr, ok := balance.(string)
					assert.True(t, ok, "Expected balance to be a string")

					// Check that the balance is a valid number (including zero)
					// We're using big.Int to handle potentially large balance values
					bigBalance, success := new(big.Int).SetString(balanceStr, 10)
					assert.True(t, success, "Expected balance to be a valid number")
					assert.NotNil(t, bigBalance, "Expected non-nil balance")

					// Check that the balance is non-negative
					assert.True(t, bigBalance.Sign() >= 0, "Expected non-negative balance")
				}
			},
			ConfigForResolver: cfgResolverWithOneInchApi(),
		},
		// GetWalletBalanceForUniswapV3Factory: This test case retrieves the wallet balance for the Uniswap V3 Factory address. It ensures the result is a non-empty hex string starting with "0x", using the Infura API.
		{
			Name:      "GetWalletBalanceForUniswapV3Factory",
			SessionId: generateSessionID(),
			JsonRequest: &types.JsonRequest{
				Id:     generateRequestID(),
				Method: "GetWalletBalance",
				Params: []string{"1", "0x1F98431c8aD98523631AE4a59f267346ea31F984"}, // Uniswap V3 Factory
			},
			FuncCheckActualJsonResponseResult: func(result interface{}) {
				assert.NotEmpty(t, result)
				// For OneInch, result should be a map of token addresses to balances
				balanceMap, ok := result.(map[string]interface{})
				assert.True(t, ok, "Expected result to be a map")

				// Check that the map is not empty
				assert.NotEmpty(t, balanceMap, "Expected non-empty balance map")

				// Check each key-value pair in the map
				for tokenAddress, balance := range balanceMap {
					// Check that the token address is a valid Ethereum address
					assert.Regexp(t, "^0x[a-fA-F0-9]{40}$", tokenAddress, "Expected token address to be a valid Ethereum address")

					// Check that the balance is a string
					balanceStr, ok := balance.(string)
					assert.True(t, ok, "Expected balance to be a string")

					// Check that the balance is a valid number (including zero)
					// We're using big.Int to handle potentially large balance values
					bigBalance, success := new(big.Int).SetString(balanceStr, 10)
					assert.True(t, success, "Expected balance to be a valid number")
					assert.NotNil(t, bigBalance, "Expected non-nil balance")

					// Check that the balance is non-negative
					assert.True(t, bigBalance.Sign() >= 0, "Expected non-negative balance")
				}
			},
			ConfigForResolver: cfgResolverWithOneInchApi(),
		},
		{
			Name:      "ResolverUsedDefaultHandler",
			SessionId: generateSessionID(),
			JsonRequest: &types.JsonRequest{
				Id:     generateRequestID(),
				Method: "GetWalletBalance",
				Params: []string{"0x0ADfCCa4B2a1132F82488546AcA086D7E24EA324", "latest"},
			},
			FuncCheckActualJsonResponseResult: func(result interface{}) {
				assert.Equal(t, 555., result)
			},
			ConfigForResolver: cfgResolverWithDefaultApi(),
		},
		// ResolverUsedDefaultHandler_ZeroAddress: This test case checks the wallet balance retrieval for the zero address using the default API. It verifies that the result is equal to 555.
		{
			Name:      "ResolverUsedDefaultHandler_ZeroAddress",
			SessionId: generateSessionID(),
			JsonRequest: &types.JsonRequest{
				Id:     generateRequestID(),
				Method: "GetWalletBalance",
				Params: []string{"0x0000000000000000000000000000000000000000", "latest"},
			},
			FuncCheckActualJsonResponseResult: func(result interface{}) {
				assert.Equal(t, 555., result)
			},
			ConfigForResolver: cfgResolverWithDefaultApi(),
		},
		// ResolverUsedDefaultHandler_EarliestBlock: This test case retrieves the wallet balance for a specific address at the "earliest" block state using the default API. It ensures the result is equal to 555.
		{
			Name:      "ResolverUsedDefaultHandler_EarliestBlock",
			SessionId: generateSessionID(),
			JsonRequest: &types.JsonRequest{
				Id:     generateRequestID(),
				Method: "GetWalletBalance",
				Params: []string{"0x0ADfCCa4B2a1132F82488546AcA086D7E24EA324", "earliest"},
			},
			FuncCheckActualJsonResponseResult: func(result interface{}) {
				assert.Equal(t, 555., result)
			},
			ConfigForResolver: cfgResolverWithDefaultApi(),
		},
		// ResolverUsedDefaultHandler_PendingBlock: This test case fetches the wallet balance for a specific address at the "pending" block state using the default API. It confirms that the result is equal to 555.
		{
			Name:      "ResolverUsedDefaultHandler_PendingBlock",
			SessionId: generateSessionID(),
			JsonRequest: &types.JsonRequest{
				Id:     generateRequestID(),
				Method: "GetWalletBalance",
				Params: []string{"0x0ADfCCa4B2a1132F82488546AcA086D7E24EA324", "pending"},
			},
			FuncCheckActualJsonResponseResult: func(result interface{}) {
				assert.Equal(t, 555., result)
			},
			ConfigForResolver: cfgResolverWithDefaultApi(),
		},
		// ResolverUsedDefaultHandler_SpecificBlockNumber: This test case retrieves the wallet balance for a specific address at a given block number (1,000,000) using the default API. It checks that the result is equal to 555.
		{
			Name:      "ResolverUsedDefaultHandler_SpecificBlockNumber",
			SessionId: generateSessionID(),
			JsonRequest: &types.JsonRequest{
				Id:     generateRequestID(),
				Method: "GetWalletBalance",
				Params: []string{"0x0ADfCCa4B2a1132F82488546AcA086D7E24EA324", "1000000"},
			},
			FuncCheckActualJsonResponseResult: func(result interface{}) {
				assert.Equal(t, 555., result)
			},
			ConfigForResolver: cfgResolverWithDefaultApi(),
		},
		// ResolverUsedDefaultHandler_InvalidAddress: This test case checks the wallet balance retrieval for an invalid address using the default API. It verifies that the result is equal to 555.
		{
			Name:      "ResolverUsedDefaultHandler_InvalidAddress",
			SessionId: generateSessionID(),
			JsonRequest: &types.JsonRequest{
				Id:     generateRequestID(),
				Method: "GetWalletBalance",
				Params: []string{"0xInvalidAddress", "latest"},
			},
			FuncCheckActualJsonResponseResult: func(result interface{}) {
				assert.Equal(t, 555., result)
			},
			ConfigForResolver: cfgResolverWithDefaultApi(),
		},
		{
			Name:      "ResolverUsedInfuraHandler",
			SessionId: generateSessionID(),
			JsonRequest: &types.JsonRequest{
				Id:     generateRequestID(),
				Method: "GetWalletBalance",
				Params: []string{"0x0ADfCCa4B2a1132F82488546AcA086D7E24EA324", "latest"},
			},
			FuncCheckActualJsonResponseResult: func(result interface{}) {
				assert.NotEmpty(t, result)
			},
			ConfigForResolver: cfgResolverWithInfuraApi(),
		},
		// GetWalletBalanceForVitalikButerinAddress: This test case checks the wallet balance retrieval for Vitalik Buterin's address using the Infura API. It verifies that the result is a non-empty hex string starting with "0x".
		{
			Name:      "GetWalletBalanceForVitalikButerinAddress",
			SessionId: generateSessionID(),
			JsonRequest: &types.JsonRequest{
				Id:     generateRequestID(),
				Method: "GetWalletBalance",
				Params: []string{"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", "latest"}, // Vitalik Buterin's address
			},
			FuncCheckActualJsonResponseResult: func(result interface{}) {
				assert.NotEmpty(t, result)
				// For Infura, result should be a hex string
				hexStr, ok := result.(string)
				assert.True(t, ok, "Expected balance to be a string")
				assert.True(t, len(hexStr) > 2 && hexStr[:2] == "0x", "Expected hex string starting with 0x")
			},
			ConfigForResolver: cfgResolverWithInfuraApi(),
		},
		// GetWalletBalanceForUniswapV3Factory: This test case retrieves the wallet balance for the Uniswap V3 Factory address. It ensures the result is a non-empty hex string starting with "0x", using the Infura API.
		{
			Name:      "GetWalletBalanceForUniswapV3Factory",
			SessionId: generateSessionID(),
			JsonRequest: &types.JsonRequest{
				Id:     generateRequestID(),
				Method: "GetWalletBalance",
				Params: []string{"0x1F98431c8aD98523631AE4a59f267346ea31F984", "latest"}, // Uniswap V3 Factory
			},
			FuncCheckActualJsonResponseResult: func(result interface{}) {
				assert.NotEmpty(t, result)
				// For Infura, result should be a hex string
				hexStr, ok := result.(string)
				assert.True(t, ok, "Expected balance to be a string")
				assert.True(t, len(hexStr) > 2 && hexStr[:2] == "0x", "Expected hex string starting with 0x")
			},
			ConfigForResolver: cfgResolverWithInfuraApi(),
		},
		// GetWalletBalanceForEthereumFoundation: This test case fetches the wallet balance for the Ethereum Foundation address. It confirms that the result is a non-empty hex string starting with "0x", using the Infura API.
		{
			Name:      "GetWalletBalanceForEthereumFoundation",
			SessionId: generateSessionID(),
			JsonRequest: &types.JsonRequest{
				Id:     generateRequestID(),
				Method: "GetWalletBalance",
				Params: []string{"0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe", "latest"}, // Ethereum Foundation
			},
			FuncCheckActualJsonResponseResult: func(result interface{}) {
				assert.NotEmpty(t, result)
				// For Infura, result should be a hex string
				hexStr, ok := result.(string)
				assert.True(t, ok, "Expected balance to be a string")
				assert.True(t, len(hexStr) > 2 && hexStr[:2] == "0x", "Expected hex string starting with 0x")
			},
			ConfigForResolver: cfgResolverWithInfuraApi(),
		},
		// GetWalletBalanceWithPendingBlock: This test case retrieves the wallet balance for a specific address at the "pending" block state. It checks that the result is a non-empty hex string starting with "0x", using the Infura API.
		{
			Name:      "GetWalletBalanceWithPendingBlock",
			SessionId: generateSessionID(),
			JsonRequest: &types.JsonRequest{
				Id:     generateRequestID(),
				Method: "GetWalletBalance",
				Params: []string{"0x0ADfCCa4B2a1132F82488546AcA086D7E24EA324", "pending"},
			},
			FuncCheckActualJsonResponseResult: func(result interface{}) {
				assert.NotEmpty(t, result)
				// For Infura, result should be a hex string
				hexStr, ok := result.(string)
				assert.True(t, ok, "Expected balance to be a string")
				assert.True(t, len(hexStr) > 2 && hexStr[:2] == "0x", "Expected hex string starting with 0x")
			},
			ConfigForResolver: cfgResolverWithInfuraApi(),
		},
		// GetWalletBalanceWithEarliestBlock: This test case fetches the wallet balance for a specific address at the "earliest" block state. It verifies that the result is a non-empty hex string starting with "0x", using the Infura API.
		{
			Name:      "GetWalletBalanceWithEarliestBlock",
			SessionId: generateSessionID(),
			JsonRequest: &types.JsonRequest{
				Id:     generateRequestID(),
				Method: "GetWalletBalance",
				Params: []string{"0x0ADfCCa4B2a1132F82488546AcA086D7E24EA324", "earliest"},
			},
			FuncCheckActualJsonResponseResult: func(result interface{}) {
				assert.NotEmpty(t, result)
				// For Infura, result should be a hex string
				hexStr, ok := result.(string)
				assert.True(t, ok, "Expected balance to be a string")
				assert.True(t, len(hexStr) > 2 && hexStr[:2] == "0x", "Expected hex string starting with 0x")
			},
			ConfigForResolver: cfgResolverWithInfuraApi(),
		},
		// GetWalletBalanceWithSpecificBlockNumber: This test case retrieves the wallet balance for a specific address at a given block number (14,680,064). It ensures the result is a non-empty hex string starting with "0x", using the Infura API.
		{
			Name:      "GetWalletBalanceWithSpecificBlockNumber",
			SessionId: generateSessionID(),
			JsonRequest: &types.JsonRequest{
				Id:     generateRequestID(),
				Method: "GetWalletBalance",
				Params: []string{"0x0ADfCCa4B2a1132F82488546AcA086D7E24EA324", "0xE10000"}, // Block 14,680,064
			},
			FuncCheckActualJsonResponseResult: func(result interface{}) {
				assert.NotEmpty(t, result)
				// For Infura, result should be a hex string
				hexStr, ok := result.(string)
				assert.True(t, ok, "Expected balance to be a string")
				assert.True(t, len(hexStr) > 2 && hexStr[:2] == "0x", "Expected hex string starting with 0x")
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

func cfgResolverWithOneInchApi() *resolver.Config {
	cfgResolver := cfgResolverWithoutApis()

	cfgResolver.Apis = resolver.ApiConfigs{
		OneInch: resolver.OneInchApiConfig{
			Enabled: true,
			Key:     getOneInchPortalTokenFromEnv(),
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
		resolverPrivKey := tn.ResolverPrivateKeys[0]
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

		privKey, err := crypto.HexToECDSA(resolverPrivKey)
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
	}, testnetwork.WithNodeRegistry())
}

func getOneInchPortalTokenFromEnv() string {
	if token, ok := os.LookupEnv("DEV_PORTAL_TOKEN"); ok {
		return token
	}

	panic("For 1inch tests need set devv portal token to environment, expected key: 'DEV_PORTAL_TOKEN'")
}

func generateUniqueID(prefix string) string {
	timestamp := time.Now().UnixNano()
	randomNum := rand.Intn(1000000)
	data := []byte(fmt.Sprintf("%s-%d-%d", prefix, timestamp, randomNum))
	hash := sha256.Sum256(data)
	return prefix + "-" + hex.EncodeToString(hash[:8])
}

func generateSessionID() string {
	return generateUniqueID("session")
}

func generateRequestID() string {
	return generateUniqueID("req")
}

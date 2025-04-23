package resolver

import (
	"log/slog"
	"os"
	"testing"

	"github.com/1inch/p2p-network/resolver/types"
	"github.com/stretchr/testify/assert"
)

const (
	getWalletBalanceMethod     = "GetWalletBalance"
	ethereumChainId            = "1"
	ethereumAddressFromMainnet = "0x4838B106FCe9647Bdf1E7877BF73cE8B0BAD5f97"
)

func TestSuccessExecuteMethod(t *testing.T) {
	logger := slog.Default()
	cfg := OneInchApiConfig{
		Key:     getOneInchPortalTokenFromEnv(),
		Enabled: true,
	}

	handler := NewOneInchApiHandler(cfg, logger)

	req := &types.JsonRequest{
		Id:     "request-id",
		Method: getWalletBalanceMethod,
		Params: []string{ethereumChainId, ethereumAddressFromMainnet},
	}
	resp, err := handler.Process(req)

	assert.Nil(t, err, "unexpected error")
	assert.NotNil(t, resp, "expected response from handler")
	assert.Equal(t, req.Id, resp.Id)
	assert.NotNil(t, resp.Result, "expected some balance")
}

func TestNegativeCases(t *testing.T) {
	logger := slog.Default()
	cfg := OneInchApiConfig{
		Key:     "incorrect-token",
		Enabled: true,
	}

	handler := NewOneInchApiHandler(cfg, logger)

	negativeTestCases := []struct {
		name        string
		request     *types.JsonRequest
		expectedErr string
	}{
		{
			name: "Unrecognized method name",
			request: &types.JsonRequest{
				Id:     "unrecognized-method name",
				Method: "unrecognized-method",
				Params: []string{ethereumChainId, ethereumAddressFromMainnet},
			},
			expectedErr: "unrecognized method",
		},
		{
			name: "Count params less than need",
			request: &types.JsonRequest{
				Id:     "count-params-less-than-need",
				Method: getWalletBalanceMethod,
				Params: []string{ethereumAddressFromMainnet},
			},
			expectedErr: "wrong number of params",
		},
		{
			name: "Address param is empty",
			request: &types.JsonRequest{
				Id:     "empty-address",
				Method: getWalletBalanceMethod,
				Params: []string{ethereumChainId, ""},
			},
			expectedErr: "empty address",
		},
		{
			name: "ChainId param is empty",
			request: &types.JsonRequest{
				Id:     "empty-chainId",
				Method: getWalletBalanceMethod,
				Params: []string{"", ethereumAddressFromMainnet},
			},
			expectedErr: "empty chainId",
		},
		{
			name: "ChainId is not numeric",
			request: &types.JsonRequest{
				Id:     "chainId-not-numeric",
				Method: getWalletBalanceMethod,
				Params: []string{"not-numeric", ethereumAddressFromMainnet},
			},
			expectedErr: "chainId must be numeric",
		},
		{
			name: "ChainId not supported",
			request: &types.JsonRequest{
				Id:     "chainId-not-supported",
				Method: getWalletBalanceMethod,
				Params: []string{"11155111", ethereumAddressFromMainnet},
			},
			expectedErr: "chainId not supported",
		},
		{
			// in handler set incorrect token, so need just put valid request
			name: "Incorrect dev portal token",
			request: &types.JsonRequest{
				Id:     "incorrect-token",
				Method: getWalletBalanceMethod,
				Params: []string{ethereumChainId, ethereumAddressFromMainnet},
			},
			// TODO change this message after receive token
			expectedErr: "processing response failed: failed to unmarshal error response body: unexpected end of JSON input",
		},
	}

	for _, testCase := range negativeTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			response, err := handler.Process(testCase.request)

			assert.NotNil(t, response, "expected response from handler")
			assert.Equal(t, testCase.request.Id, response.Id)
			assert.Equal(t, 0, response.Result)
			assert.EqualError(t, err, testCase.expectedErr, "errors not match")
		})
	}
}

func getOneInchPortalTokenFromEnv() string {
	if token, ok := os.LookupEnv("DEV_PORTAL_TOKEN"); ok {
		return token
	}

	panic("For 1inch tests need set dev portal token to environment, expected key: 'DEV_PORTAL_TOKEN'")
}

package resolver

import (
    "errors"
    "github.com/stretchr/testify/assert"
	"github.com/1inch/p2p-network/resolver/types"
    "log/slog"
    "os"
    "testing"
    "strconv"
    "strings"
)

type defaultApiHandler struct {
    logger *slog.Logger
}

var (
    errWrongParamCount    = errors.New("wrong number of parameters")
    errUnrecognizedMethod = errors.New("unrecognized method")
    errEmptyAddress       = errors.New("empty address")
    errEmptyChainId       = errors.New("empty chainId")
    errChainIdNotNumeric  = errors.New("chainId is not numeric")
    errInvalidParamOrder  = errors.New("invalid parameter order: expected (chainId, address)")
)

func (h *defaultApiHandler) Process(req *types.JsonRequest) (*types.JsonResponse, error) {
    if req.Method != "GetWalletBalance" {
        return nil, errUnrecognizedMethod
    }
    if len(req.Params) != 2 {
        return nil, errWrongParamCount
    }
    if req.Params[0] == "" {
        return nil, errEmptyChainId
    }
    if req.Params[1] == "" {
        return nil, errEmptyAddress
    }
    if strings.HasPrefix(req.Params[0], "0x") {
        return nil, errInvalidParamOrder
    }
    if !isNumeric(req.Params[0]) {
        return nil, errChainIdNotNumeric
    }

    return &types.JsonResponse{
        Id:     req.Id,
        Result: 555,
    }, nil
}

// The test cases cover:
// - Validate correct handling of parameter order, ensuring the chainId comes before the address.
// - Check that unrecognized method names are properly rejected.
// - Verify that valid parameters result in the expected response.
// - Test handling of incorrect parameter counts, both too few and too many.
// - Ensure proper handling of empty or nil parameter lists.
// - Verify that different valid addresses and chain IDs are accepted.
// - Test the handler's response to various block number formats.
// - Check handling of edge cases like zero addresses and very long addresses.
// - Validate proper error responses for empty chainId or address parameters.
// - Test rejection of non-numeric chainId values.
// - Verify acceptance of addresses without the '0x' prefix.
// - Check handling of very large block numbers.
// - Test rejection of non-numeric characters in the block number.
// - Verify handling of whitespace in parameters.
// - Test case sensitivity in address handling.
// RUN: go test -v ./resolver/handler_test.go
func TestDefaultApiHandler_GetWalletBalance(t *testing.T) {
    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
    handler := &defaultApiHandler{logger: logger}

    testCases := []struct {
        name           string
		request        *types.JsonRequest
        expectedResult int
        expectedError  error
    }{
        {
            name: "Incorrect params order",
            request: &types.JsonRequest{
                Method: "GetWalletBalance",
                Params: []string{"0x123", "1"},
            },
            expectedResult: 0,
            expectedError:  errInvalidParamOrder,
        },
		{
            name: "Unrecognized method name",
            request: &types.JsonRequest{
                Method: "UnknownMethod",
                Params: []string{"1", "0x123"},
            },
            expectedResult: 0,
            expectedError:  errUnrecognizedMethod,
        },
		{
            name: "Valid params",
            request: &types.JsonRequest{
                Method: "GetWalletBalance",
                Params: []string{"1", "0x123"},
            },
            expectedResult: 555,
            expectedError:  nil,
        },
        {
            name: "Wrong number of params - too few",
            request: &types.JsonRequest{
                Method: "GetWalletBalance",
                Params: []string{"0x123"},
            },
            expectedResult: 0,
            expectedError:  errWrongParamCount,
        },
        {
            name: "Wrong number of params - too many",
            request: &types.JsonRequest{
                Method: "GetWalletBalance",
                Params: []string{"0", "0x123", "extra"},
            },
            expectedResult: 0,
            expectedError:  errWrongParamCount,
        },
        {
            name: "Empty params",
            request: &types.JsonRequest{
                Method: "GetWalletBalance",
                Params: []string{},
            },
            expectedResult: 0,
            expectedError:  errWrongParamCount,
        },
        {
            name: "Nil params",
            request: &types.JsonRequest{
                Method: "GetWalletBalance",
                Params: nil,
            },
            expectedResult: 0,
            expectedError:  errWrongParamCount,
        },
        {
            name: "Valid params with different address",
            request: &types.JsonRequest{
                Method: "GetWalletBalance",
                Params: []string{"0", "0xabc"},
            },
            expectedResult: 555,
            expectedError:  nil,
        },
        {
            name: "Valid params with different block",
            request: &types.JsonRequest{
                Method: "GetWalletBalance",
                Params: []string{"1", "0x123"},
            },
            expectedResult: 555,
            expectedError:  nil,
        },
        {
            name: "Valid params with numeric block",
            request: &types.JsonRequest{
                Method: "GetWalletBalance",
                Params: []string{"123", "0x123"},
            },
            expectedResult: 555,
            expectedError:  nil,
        },
        {
            name: "Valid params with zero address",
            request: &types.JsonRequest{
                Method: "GetWalletBalance",
                Params: []string{"0", "0x0000000000000000000000000000000000000000"},
            },
            expectedResult: 555,
            expectedError:  nil,
        },
		{
            name: "Count params less than needed",
            request: &types.JsonRequest{
                Method: "GetWalletBalance",
                Params: []string{"0x123"},
            },
            expectedResult: 0,
            expectedError:  errWrongParamCount,
        },
        {
            name: "Address param is empty",
            request: &types.JsonRequest{
                Method: "GetWalletBalance",
                Params: []string{"1", ""},
            },
            expectedResult: 0,
            expectedError:  errEmptyAddress,
        },
        {
            name: "ChainId param is empty",
            request: &types.JsonRequest{
                Method: "GetWalletBalance",
                Params: []string{"", "0x123"},
            },
            expectedResult: 0,
            expectedError:  errEmptyChainId,
        },
        {
            name: "ChainId is not numeric",
            request: &types.JsonRequest{
                Method: "GetWalletBalance",
                Params: []string{"not-numeric", "0x123"},
            },
            expectedResult: 0,
            expectedError:  errChainIdNotNumeric,
        },
        {
            name: "Very long address",
            request: &types.JsonRequest{
                Method: "GetWalletBalance",
                Params: []string{"0", "0x" + strings.Repeat("1", 100)},
            },
            expectedResult: 555,
            expectedError:  nil,
        },
        {
            name: "Address without 0x prefix",
            request: &types.JsonRequest{
                Method: "GetWalletBalance",
                Params: []string{"0", "123abc"},
            },
            expectedResult: 555,
            expectedError:  nil,
        },
        {
            name: "Very large block number",
            request: &types.JsonRequest{
                Method: "GetWalletBalance",
                Params: []string{"999999999999999", "0x123"},
            },
            expectedResult: 555,
            expectedError:  nil,
        },
        {
            name: "Non-numeric characters in block number",
            request: &types.JsonRequest{
                Method: "GetWalletBalance",
                Params: []string{"12a34", "0x123"},
            },
            expectedResult: 0,
            expectedError:  errChainIdNotNumeric,
        },
        {
            name: "Whitespace in parameters",
            request: &types.JsonRequest{
                Method: "GetWalletBalance",
                Params: []string{" 0 ", " 0x123 "},
            },
            expectedResult: 555,
            expectedError:  errChainIdNotNumeric,
        },
        {
            name: "Case sensitivity test",
            request: &types.JsonRequest{
                Method: "GetWalletBalance",
                Params: []string{"0", "0xAbC123"},
            },
            expectedResult: 555,
            expectedError:  nil,
        },
    }

	for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result, err := handler.Process(tc.request)

            if tc.expectedError != nil {
                assert.Error(t, err)
                assert.Equal(t, tc.expectedError, err)
                assert.Nil(t, result)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, result)
                assert.Equal(t, tc.expectedResult, result.Result)
            }
        })
    }
}

func isNumeric(s string) bool {
    if s == "" {
        return false
    }
    _, err := strconv.ParseInt(s, 10, 64)
    return err == nil
}
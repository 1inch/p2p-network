package resolver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/1inch/p2p-network/resolver/types"
)

var errUnexpectedStatusCode = errors.New("unexpected status code")
var errUnauthorizedStatusCode = errors.New("unauthorized call api")

const (
	baseUrl                         = "https://api.1inch.dev/"
	routingFormatToGetWalletBalance = "balance/v1.2/%s/balances/%s" // 1 - chainId, 2 - walletAddress
)

type oneInchApiHandler struct {
	cfg    OneInchApiConfig
	client *http.Client
	logger *slog.Logger
}

// NewOneInchApiHandler creates an 1inch API handler instance
func NewOneInchApiHandler(cfg OneInchApiConfig, logger *slog.Logger) ApiHandler {
	return &oneInchApiHandler{
		cfg:    cfg,
		logger: logger.WithGroup("1inch-handler-api"),
		client: http.DefaultClient,
	}
}

// Process acts as an API wrapper for JSON payloads coming through gRPC
func (h *oneInchApiHandler) Process(req *types.JsonRequest) (*types.JsonResponse, error) {
	switch req.Method {
	case "GetWalletBalance":
		balance, err := h.getWalletBalance(req.Params)
		if err != nil {
			return &types.JsonResponse{Id: req.Id, Result: 0}, err
		}
		return &types.JsonResponse{Id: req.Id, Result: balance}, nil
	default:
		return &types.JsonResponse{Id: req.Id, Result: 0}, errUnrecognizedMethod
	}
}

func (h *oneInchApiHandler) getWalletBalance(params []string) (string, error) {
	if len(params) != 2 {
		h.logger.Error("GetWalletBalance: wrong number of params", "len", len(params))
		return "", errWrongParamCount
	}
	chainId := params[0]
	address := params[1]

	err := h.validateRequest(chainId, address)
	if err != nil {
		h.logger.Error("failed validate request for GetWalletBalance")
		return "", err
	}
	request, err := h.formatHttpRequest(chainId, address)
	if err != nil {
		return "", err
	}

	resp, err := h.client.Do(request)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			h.logger.Error("failed close response bodys")
		}
	}()
	if err != nil {
		h.logger.Error("failed invoking JSON-RPC request", slog.Any("err", err.Error()))
		return "", err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		h.logger.Error("problems with authorization token, returned unauthorized status code")
		return "", errUnauthorizedStatusCode
	}

	if resp.StatusCode != http.StatusOK {
		h.logger.Error("returned invalid status code from 1inch", slog.Any("status", resp.Status))
		return "", errUnexpectedStatusCode
	}

	return h.decodeBalance(resp)
}

func (h *oneInchApiHandler) validateRequest(chainId, address string) error {
	if address == "" {
		return errEmptyAddress
	}

	if chainId == "" {
		return errEmptyChainId
	}

	return nil
}

func (h *oneInchApiHandler) formatHttpRequest(chainId, address string) (*http.Request, error) {
	requestUrl := baseUrl + fmt.Sprintf(routingFormatToGetWalletBalance, chainId, address)
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, requestUrl, nil)

	if err != nil {
		h.logger.Error("failed create http request", slog.Any("err", err.Error()))
		return nil, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.cfg.Key))
	request.Header.Add("Content-Type", "application/json")

	return request, nil
}

func (h *oneInchApiHandler) decodeBalance(resp *http.Response) (string, error) {
	var balance string

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&balance)

	if err != nil {
		h.logger.Error("failed decode wallet balance", slog.Any("err", err.Error()))
		return "", err
	}

	h.logger.Info("GetWalletBalance() processed")
	return balance, nil
}

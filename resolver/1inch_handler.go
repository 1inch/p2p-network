package resolver

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/1inch/1inch-sdk-go/sdk-clients/balances"
	"github.com/1inch/p2p-network/resolver/types"
)

const apiUrl = "https://api.1inch.dev"

// oneInchApiHandler represends handler which would call 1inch api
type oneInchApiHandler struct {
	cfg            OneInchApiConfig
	logger         *slog.Logger
	balanceClients map[string]*balances.Client // map<chainId, balances.Client>
}

// NewOneInchApiHandler creates an 1inch API handler instance
func NewOneInchApiHandler(cfg OneInchApiConfig, logger *slog.Logger) ApiHandler {
	return &oneInchApiHandler{
		cfg:    cfg,
		logger: logger.WithGroup("1inch-handler-api"),
	}
}

func (h *oneInchApiHandler) getClientByChainId(chainId string) (*balances.Client, error) {
	if client, ok := h.balanceClients[chainId]; ok {
		return client, nil
	}

	chainIdInt, err := strconv.ParseUint(chainId, 10, 64)
	if err != nil {
		h.logger.Error("Invalid chainId format", slog.Any("chainId", chainId))
		return nil, err
	}

	config, err := balances.NewConfiguration(
		balances.ConfigurationParams{
			ChainId: chainIdInt,
			ApiUrl:  apiUrl,
			ApiKey:  h.cfg.Key,
		},
	)

	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed create configuration for client of balances for chainId: %s", chainId), slog.Any("err", err.Error()))
		return nil, err
	}

	newClient, err := balances.NewClient(config)

	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed create client of balances for chainId: %s", chainId), slog.Any("err", err.Error()))
		return nil, err
	}

	h.balanceClients[chainId] = newClient

	return newClient, nil
}

// Process acts as an API wrapper for JSON payloads coming through gRPC
func (h *oneInchApiHandler) Process(req *types.JsonRequest) (*types.JsonResponse, error) {
	switch req.Method {
	case "GetWalletBalance":
		balances, err := h.getWalletBalance(req.Params)
		if err != nil {
			return &types.JsonResponse{Id: req.Id, Result: 0}, err
		}
		return &types.JsonResponse{Id: req.Id, Result: balances}, nil
	default:
		return &types.JsonResponse{Id: req.Id, Result: 0}, errUnrecognizedMethod
	}
}

func (h *oneInchApiHandler) getWalletBalance(params []string) (interface{}, error) {
	if len(params) != 2 {
		h.logger.Error("GetWalletBalance: wrong number of params", "len", len(params))
		return nil, errWrongParamCount
	}
	chainId := params[0]
	address := params[1]

	err := h.validateRequest(chainId, address)
	if err != nil {
		h.logger.Error("failed validate request for GetWalletBalance")
		return nil, err
	}

	client, err := h.getClientByChainId(chainId)

	if err != nil {
		return nil, err
	}

	resp, err := client.GetBalancesByWalletAddress(
		context.Background(),
		balances.BalancesByWalletAddressParams{
			Wallet: address,
		},
	)

	if err != nil {
		h.logger.Error("failed invoking JSON-RPC request", slog.Any("err", err.Error()))
		return nil, err
	}

	return resp, nil
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

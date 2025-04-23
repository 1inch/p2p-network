package resolver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/1inch/1inch-sdk-go/sdk-clients/balances"
	"github.com/1inch/p2p-network/resolver/types"
)

const (
	apiUrl          = "https://api.1inch.dev"
	chainIdArbitrum = "42161"
	chainIdAurora   = "1313161554"
	chainIdAvalance = "43114"
	chainIdBase     = "8453"
	chainIdBinance  = "56"
	chainIdZkSync   = "324"
	chainIdEthereum = "1"
	chainIdFantom   = "250"
	chainIdGnosis   = "100"
	chainIdKaia     = "8217"
	chainOptimism   = "10"
	chainIdPolygon  = "137"
	chainIdLinea    = "59144"
)

var (
	errChainIdMustBeNumeric = errors.New("chainId must be numeric")
	errChainIdNotSupported  = errors.New("chainId not supported")
)

// oneInchApiHandler represends handler which would call 1inch api
type oneInchApiHandler struct {
	cfg            OneInchApiConfig
	logger         *slog.Logger
	balanceClients map[uint64]*balances.Client // map<chainId, balances.Client>
}

// NewOneInchApiHandler creates an 1inch API handler instance
func NewOneInchApiHandler(cfg OneInchApiConfig, logger *slog.Logger) ApiHandler {
	return &oneInchApiHandler{
		cfg:            cfg,
		logger:         logger.WithGroup("1inch-handler-api"),
		balanceClients: make(map[uint64]*balances.Client),
	}
}

func (h *oneInchApiHandler) getClientByChainId(chainId uint64) (*balances.Client, error) {
	if client, ok := h.balanceClients[chainId]; ok {
		return client, nil
	}

	config, err := balances.NewConfiguration(
		balances.ConfigurationParams{
			ChainId: chainId,
			ApiUrl:  apiUrl,
			ApiKey:  h.cfg.Key,
		},
	)

	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed create configuration for client of balances for chainId: %d", chainId), slog.Any("err", err.Error()))
		return nil, err
	}

	newClient, err := balances.NewClient(config)

	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed create client of balances for chainId: %d", chainId), slog.Any("err", err.Error()))
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

	chainIdInt, err := h.validateRequest(chainId, address)
	if err != nil {
		h.logger.Error("failed validate request for GetWalletBalance")
		return nil, err
	}

	client, err := h.getClientByChainId(chainIdInt)

	if err != nil {
		return nil, err
	}

	// TODO add handler for errors with mapping after receive token
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

func (h *oneInchApiHandler) validateRequest(chainId, address string) (uint64, error) {
	if address == "" {
		return 0, errEmptyAddress
	}

	if chainId == "" {
		return 0, errEmptyChainId
	}

	chainIdInt, err := strconv.ParseUint(chainId, 10, 64)
	if err != nil {
		h.logger.Error("Invalid chainId format", slog.Any("chainId", chainId))
		return 0, errChainIdMustBeNumeric
	}

	if !h.isSupportedChainId(chainId) {
		return 0, errChainIdNotSupported
	}

	return chainIdInt, nil
}

func (h *oneInchApiHandler) isSupportedChainId(chainId string) bool {
	switch chainId {
	case chainIdArbitrum, chainIdAvalance, chainIdAurora, chainIdBase, chainIdBinance,
		chainIdEthereum, chainIdFantom, chainIdGnosis, chainIdKaia, chainIdLinea, chainIdPolygon,
		chainIdZkSync, chainOptimism:
		{
			return true
		}
	default:
		return false
	}
}

package resolver

import (
	"log/slog"

	"github.com/1inch/p2p-network/resolver/types"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

type infuraApiHandler struct {
	client *gethrpc.Client
	logger *slog.Logger
}

// NewInfuraApiHandler creates an Infura API handler instance
func NewInfuraApiHandler(cfg InfuraApiConfig, logger *slog.Logger) ApiHandler {
	client, err := gethrpc.DialHTTP("https://mainnet.infura.io/v3/" + cfg.Key)
	if err != nil {
		slog.Error("error creating Infura API client", "err", err)
		return nil
	}
	return &infuraApiHandler{client, logger.With("module", "api-infura")}
}

// Process acts as an API wrapper for JSON payloads coming through gRPC
func (h *infuraApiHandler) Process(req *types.JsonRequest) (*types.JsonResponse, error) {
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

func (h *infuraApiHandler) getWalletBalance(params []string) (string, error) {
	if len(params) != 2 {
		h.logger.Error("GetWalletBalance: wrong number of params", "cnt", len(params))
		return "", errWrongParamCount
	}
	address := params[0]
	block := params[1]

	err := h.validateRequest(address, block)
	if err != nil {
		h.logger.Error("failed validate request for GetWalletBalance")
		return "", err
	}

	var result string
	err = h.client.Call(&result, "eth_getBalance", address, block)
	if err != nil {
		h.logger.Error("failed invoking JSON-RPC request", "err", err)
		return "", err
	}
	h.logger.Info("GetWalletBalance() processed")
	return result, nil
}

func (h *infuraApiHandler) validateRequest(address, block string) error {
	if address == "" {
		return errEmptyAddress
	}

	if block == "" {
		return errEmptyBlock
	}

	return nil
}

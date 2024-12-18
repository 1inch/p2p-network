package resolver

import (
	"log/slog"
	"os"

	"github.com/1inch/p2p-network/resolver/types"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

type infuraApiHandler struct {
	client *gethrpc.Client
	logger *slog.Logger
}

// NewInfuraApiHandler creates an Infura API handler instance
func NewInfuraApiHandler(cfg InfuraApiConfig) ApiHandler {
	client, err := gethrpc.DialHTTP("https://mainnet.infura.io/v3/" + cfg.Key)
	if err != nil {
		slog.Error("Error creating Infura API client", "err", err)
		return nil
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With("module", "api-infura")
	return &infuraApiHandler{client, logger}
}

// Name returns API name
func (h *infuraApiHandler) Name() string {
	return "infura"
}

// Process acts as an API wrapper for JSON payloads coming through gRPC
func (h *infuraApiHandler) Process(req *types.JsonRequest) *types.JsonResponse {
	switch req.Method {
	case "GetWalletBalance":
		balance, err := h.getWalletBalance(req.Params)
		if err != nil {
			return &types.JsonResponse{Id: req.Id, Result: 0, Error: err.Error()}
		}
		return &types.JsonResponse{Id: req.Id, Result: balance}
	default:
		return &types.JsonResponse{Id: req.Id, Result: 0, Error: errUnrecognizedMethod.Error()}
	}
}

func (h *infuraApiHandler) getWalletBalance(params []string) (string, error) {
	if len(params) != 2 {
		h.logger.Error("GetWalletBalance: wrong number of params", "cnt", len(params))
		return "", errWrongParamCount
	}
	address := params[0]
	block := params[1]
	var result string
	err := h.client.Call(&result, "eth_getBalance", address, block)
	if err != nil {
		h.logger.Error("Error invoking JSON-RPC request", "err", err)
		return "", err
	}
	h.logger.Info("GetWalletBalance() processed")
	return result, nil
}

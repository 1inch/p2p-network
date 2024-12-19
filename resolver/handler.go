package resolver

import (
	"errors"
	"log/slog"
	"os"

	"github.com/1inch/p2p-network/resolver/types"
)

// ApiHandler provides Process() method for handling JSON payloads
type ApiHandler interface {
	Process(*types.JsonRequest) *types.JsonResponse
	Name() string
}

type defaultApiHandler struct {
	logger *slog.Logger
}

// Name returns API name
func (h *defaultApiHandler) Name() string {
	return "default"
}

// NewDefaultApiHandler creates a default API handler instance
func NewDefaultApiHandler(cfg DefaultApiConfig) ApiHandler {
	return &defaultApiHandler{logger: slog.New(slog.NewTextHandler(os.Stdout, nil)).With("module", "api")}
}

var errUnrecognizedMethod = errors.New("Unrecognized method")
var errWrongParamCount = errors.New("Wrong number of params")

// Process acts as an API wrapper for JSON payloads coming through gRPC
func (h *defaultApiHandler) Process(req *types.JsonRequest) *types.JsonResponse {
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

func (h *defaultApiHandler) getWalletBalance(params []string) (int, error) {
	if len(params) != 2 {
		h.logger.Error("GetWalletBalance: wrong number of params", "cnt", len(params))
		return 0, errWrongParamCount
	}
	h.logger.Info("GetWalletBalance() processed")
	return 555, nil
}

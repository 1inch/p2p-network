package resolver

import (
	"errors"
	"log/slog"

	"github.com/1inch/p2p-network/resolver/types"
)

var (
	errUnrecognizedMethod = errors.New("unrecognized method")
	errWrongParamCount    = errors.New("wrong number of params")
	errEmptyAddress       = errors.New("empty address")
	errEmptyBlock         = errors.New("empty block")
)

// ApiHandler provides Process() method for handling JSON payloads
type ApiHandler interface {
	Process(*types.JsonRequest) (*types.JsonResponse, error)
}

type defaultApiHandler struct {
	logger *slog.Logger
}

// NewDefaultApiHandler creates a default API handler instance
func NewDefaultApiHandler(cfg DefaultApiConfig, logger *slog.Logger) ApiHandler {
	return &defaultApiHandler{logger: logger.With("module", "api")}
}

// Process acts as an API wrapper for JSON payloads coming through gRPC
func (h *defaultApiHandler) Process(req *types.JsonRequest) (*types.JsonResponse, error) {
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

func (h *defaultApiHandler) getWalletBalance(params []string) (int, error) {
	if len(params) != 2 {
		h.logger.Error("GetWalletBalance: wrong number of params", "cnt", len(params))
		return 0, errWrongParamCount
	}
	h.logger.Info("GetWalletBalance() processed")
	return 555, nil
}

package api

import (
	"errors"

	"github.com/1inch/p2p-network/resolver/types"
)

type Handler interface {
	Process(*types.JsonRequest) *types.JsonResponse
}

type DefaultHandler struct{}

func (h *DefaultHandler) Process(req *types.JsonRequest) *types.JsonResponse {
	switch req.Method {
	case "GetWalletBalance":
		balance, err := GetWalletBalance(req.Params)
		if err != nil {
			return &types.JsonResponse{Id: req.Id, Result: 0, Error: err}
		}
		return &types.JsonResponse{Id: req.Id, Result: balance, Error: nil}
	default:
		return &types.JsonResponse{Id: req.Id, Result: 0, Error: errors.New("Unrecognized method")}
	}
}

func GetWalletBalance(params []string) (int, error) {
	return 0, nil
}

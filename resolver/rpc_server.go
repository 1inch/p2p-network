// Package resolver implements the gRPC server
package resolver

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"log/slog"
	"os"

	pb "github.com/1inch/p2p-network/proto"
	"github.com/1inch/p2p-network/resolver/types"
	"google.golang.org/grpc"
)

// Server represents gRPC server.
type Server struct {
	pb.UnimplementedExecuteServer

	privateKey *rsa.PrivateKey

	logger *slog.Logger

	grpcServer *grpc.Server

	handlers []ApiHandler
}

func generateKey() (*rsa.PrivateKey, error) {
	p, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	return p, nil
}

var errNoHandlers = errors.New("No API handlers passed to server")

// newServer creates new RpcServer.
func newServer(handlers []ApiHandler) (*Server, error) {
	if len(handlers) == 0 {
		return nil, errNoHandlers
	}
	privKey, err := generateKey()
	if err != nil {
		return nil, err
	}
	return &Server{privateKey: privKey, logger: slog.New(slog.NewTextHandler(os.Stdout, nil)).With("module", "server"), handlers: handlers}, nil
}

func (s *Server) processRequest(req *pb.ResolverRequest) ([]byte, error) {
	// Unmarshal JSON
	var jsonReq types.JsonRequest
	err := json.Unmarshal(req.Payload, &jsonReq)
	if err != nil {
		return nil, err
	}

	jsonResponses := make(map[string]interface{})
	for _, h := range s.handlers {
		jsonResp := h.Process(&jsonReq)
		jsonResponses[h.Name()] = jsonResp
	}
	byteArr, err := json.Marshal(jsonResponses)
	if err != nil {
		return nil, err
	}

	return byteArr, nil
}

// Execute executes ResolverRequest.
func (s *Server) Execute(ctx context.Context, req *pb.ResolverRequest) (*pb.ResolverResponse, error) {
	s.logger.Info("###Incoming request", "id", req.Id)
	resp, err := s.processRequest(req)
	var respStatus pb.ResolverResponseStatus
	if err != nil {
		respStatus = pb.ResolverResponseStatus_RESOLVER_ERROR
		slog.Error("processRequest() error", "err", err)
	} else {
		respStatus = pb.ResolverResponseStatus_RESOLVER_OK
	}

	response := &pb.ResolverResponse{
		Id:      req.Id,
		Payload: resp,
		Status:  respStatus,
	}
	return response, nil
}

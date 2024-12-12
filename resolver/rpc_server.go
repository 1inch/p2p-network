// Package resolver implements the gRPC server
package resolver

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"log/slog"
	"os"

	pb "github.com/1inch/p2p-network/proto"
	"github.com/1inch/p2p-network/resolver/types"
	"google.golang.org/grpc"
)

// Config represents basic server config
type Config struct {
	Port     int
	LogLevel slog.Level
}

// Server represents gRPC server.
type Server struct {
	pb.UnimplementedExecuteServer

	privateKey *rsa.PrivateKey

	logger *slog.Logger

	grpcServer *grpc.Server

	handler ApiHandler
}

func generateKey() (*rsa.PrivateKey, error) {
	p, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// newServer creates new RpcServer.
func newServer(apiHandler ApiHandler) (*Server, error) {
	privKey, err := generateKey()
	if err != nil {
		return nil, err
	}

	return &Server{privateKey: privKey, logger: slog.New(slog.NewTextHandler(os.Stdout, nil)).With("module", "server"), handler: apiHandler}, nil
}

func (s *Server) processRequest(req *pb.ResolverRequest) ([]byte, error) {
	// Unmarshal JSON
	var jsonReq types.JsonRequest
	err := json.Unmarshal(req.Payload, &jsonReq)
	if err != nil {
		return nil, err
	}

	jsonResp := s.handler.Process(&jsonReq)
	byteArr, err := json.Marshal(jsonResp)
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

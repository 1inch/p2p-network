// Package rpc represents GRPC server.
package rpc

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"log/slog"
	"os"

	pb "github.com/1inch/p2p-network/proto"
	"google.golang.org/grpc"
)

type Config struct {
	Port     int
	LogLevel slog.Level
	Testing  bool
}

// Server represents gRPC server.
type Server struct {
	pb.UnimplementedExecuteServer

	privateKey *rsa.PrivateKey

	logger *slog.Logger

	grpcServer *grpc.Server
}

func generateKey() (*rsa.PrivateKey, error) {
	p, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// NewRpcServer creates new RpcServer.
func newServer(logLevel slog.Level) (*Server, error) {
	privKey, err := generateKey()
	if err != nil {
		return nil, err
	}

	return &Server{privateKey: privKey, logger: slog.New(slog.NewTextHandler(os.Stdout, nil))}, nil
}

// Execute executes ResolverRequest.
func (s *Server) Execute(ctx context.Context, req *pb.ResolverRequest) (*pb.ResolverResponse, error) {
	s.logger.Info("###Incoming request", "id", req.Id)
	response := &pb.ResolverResponse{
		Id:      req.Id,
		Payload: req.Payload,
		Status:  pb.ResolverResponseStatus_RESOLVER_OK,
	}
	return response, nil
}

func (s *Server) GetPublicKey(ctx context.Context) ([]byte, error) {
	byteArr, err := x509.MarshalPKIXPublicKey(s.privateKey.Public())
	if err != nil {
		return nil, err
	}
	return byteArr, nil
}

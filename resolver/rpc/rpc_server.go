// Package rpc represents GRPC server.
package rpc

import (
	"context"

	pb "github.com/1inch/p2p-network/proto"
)

// Server represents gRPC server.
type Server struct {
	pb.UnimplementedExecuteServer
}

// NewRpcServer creates new RpcServer.
func NewRpcServer() *Server {
	return &Server{}
}

// Execute executes ResolverRequest.
func (s *Server) Execute(ctx context.Context, req *pb.ResolverRequest) (*pb.ResolverResponse, error) {
	response := &pb.ResolverResponse{
		Id:      req.Id,
		Payload: make([]byte, 0),
		Status:  pb.ResolverResponseStatus_RESOLVER_OK,
	}
	return response, nil
}

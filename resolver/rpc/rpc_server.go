package rpc

import (
	"context"

	pb "github.com/1inch/p2p-network/proto"
)

type RpcServer struct {
	pb.UnimplementedExecuteServer
}

func NewRpcServer() *RpcServer {
	return &RpcServer{}
}

func (s *RpcServer) Execute(ctx context.Context, req *pb.ResolverRequest) (*pb.ResolverResponse, error) {
	response := &pb.ResolverResponse{
		Id:      req.Id,
		Payload: make([]byte, 0),
		Status:  pb.ResolverResponseStatus_RESOLVER_OK,
	}
	return response, nil
}

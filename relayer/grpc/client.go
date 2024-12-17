// Package grpc defines grpc client.
package grpc

import (
	"context"

	pb "github.com/1inch/p2p-network/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps the gRPC connection and Execute service client.
type Client struct {
	Conn          *grpc.ClientConn
	ExecuteClient pb.ExecuteClient
}

// New initializes a new gRPC client with Execute service.
func New(address string) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		Conn:          conn,
		ExecuteClient: pb.NewExecuteClient(conn),
	}, nil
}

// Close closes the gRPC connection.
func (c *Client) Close() error {
	return c.Conn.Close()
}

// Execute wraps the Execute RPC call.
func (c *Client) Execute(ctx context.Context, req *pb.ResolverRequest) (*pb.ResolverResponse, error) {
	return c.ExecuteClient.Execute(ctx, req)
}

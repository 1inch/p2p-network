// Package grpc defines grpc client.
package grpc

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"

	"github.com/1inch/p2p-network/internal/registry"
	pb "github.com/1inch/p2p-network/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// ErrResolverLookupFailed is returned when the registry client fails to resolve a public key.
	ErrResolverLookupFailed = errors.New("resolver lookup failed")
	// ErrGRPCExecutionFailed is returned when the grpc execution fails.
	ErrGRPCExecutionFailed = errors.New("gRPC execution failed")
	// ErrGRPCConnectionCloseFailed is returned when the grpc execution fails.
	ErrGRPCConnectionCloseFailed = errors.New("gRPC connection close failed")
)

// Client wraps the gRPC connection and Execute service client.
type Client struct {
	conns          map[string]*grpc.ClientConn
	registryClient *registry.Client
	mu             sync.Mutex
}

// New initializes a new gRPC client with Execute service.
func New(registryClient *registry.Client) *Client {
	return &Client{
		conns:          make(map[string]*grpc.ClientConn),
		registryClient: registryClient,
	}
}

// Execute wraps the Execute RPC call.
func (c *Client) Execute(ctx context.Context, publicKey []byte, req *pb.ResolverRequest) (*pb.ResolverResponse, error) {
	conn, err := c.getConn(publicKey)
	if err != nil {
		return nil, err
	}

	client := pb.NewExecuteClient(conn)
	response, err := client.Execute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%w: publicKey %s: %w", ErrGRPCExecutionFailed, hex.EncodeToString(publicKey), err)
	}

	return response, nil
}

// Close closes the gRPC connection.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var closeErrs []error
	for _, conn := range c.conns {
		if err := conn.Close(); err != nil {
			closeErrs = append(closeErrs, fmt.Errorf("failed to close gRPC connection: %w", err))
		}
	}

	if len(closeErrs) > 0 {
		return fmt.Errorf("%w: multiple errors: %v", ErrGRPCConnectionCloseFailed, closeErrs)
	}

	return nil
}

func (c *Client) getConn(publicKey []byte) (*grpc.ClientConn, error) {
	address, err := c.registryClient.GetResolver(publicKey)
	if err != nil {
		return nil, fmt.Errorf("%w: publicKey %s: %w", ErrResolverLookupFailed, hex.EncodeToString(publicKey), err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if conn, exists := c.conns[string(publicKey)]; exists {
		return conn, nil
	}

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c.conns[string(publicKey)] = conn
	return conn, nil
}

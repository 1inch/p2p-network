// Package grpc defines grpc client.
package grpc

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/1inch/p2p-network/internal/registry"
	pb "github.com/1inch/p2p-network/proto/resolver"
	"github.com/1inch/p2p-network/relayer/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var grpcClientConfig = `{
	"healthCheckConfig": {
		"serviceName": ""
	}
}`

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
	logger         *slog.Logger
	conns          map[string]*grpc.ClientConn
	registryClient *registry.Client
	mu             sync.Mutex
}

// New initializes a new gRPC client with Execute service.
func New(logger *slog.Logger, registryClient *registry.Client) *Client {
	return &Client{
		logger:         logger.WithGroup("grpc-server"),
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
	c.mu.Lock()
	defer c.mu.Unlock()

	if conn, exists := c.conns[string(publicKey)]; exists {
		return conn, nil
	}

	address, err := c.registryClient.GetResolver(publicKey)
	if err != nil {
		return nil, fmt.Errorf("%w: publicKey %s: %w", ErrResolverLookupFailed, hex.EncodeToString(publicKey), err)
	}

	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(grpcClientConfig),
		grpc.WithUnaryInterceptor(
			func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
				return loggingCallHandler(ctx, c.logger, method, req, reply, cc, invoker, opts...)
			}))
	if err != nil {
		return nil, err
	}

	c.conns[string(publicKey)] = conn
	return conn, nil
}

func loggingCallHandler(ctx context.Context, logger *slog.Logger, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()

	logger.Info("send request on grpc server", slog.Any("method", method))
	logger.Debug("with request", slog.Any("req", protojson.Format(req.(proto.Message))))

	err := invoker(ctx, method, req, reply, cc, opts...)
	duration := time.Since(start).Seconds()
	status := "success"
	if err != nil {
		status = "failed"
		logger.Info("request failed process", slog.Any("method", method))
		logger.Debug("with error", slog.Any("err", err.Error()))
	} else {
		logger.Info("receive response from grpc server", slog.Any("method", method))
		logger.Debug("with response", slog.Any("resp", protojson.Format(reply.(proto.Message))))
	}

	metrics.GrpcRequestsTotal.WithLabelValues(method, status).Inc()
	metrics.GrpcRequestDuration.WithLabelValues(method).Observe(duration)

	return err
}

// Package resolver represents resolver node.
package resolver

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"

	pb "github.com/1inch/p2p-network/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var errNoHandlerApiInConfig = errors.New("no handler api in config")

func setupRpcServer(listener net.Listener, server *Server, opts ...grpc.ServerOption) *grpc.Server {
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterExecuteServer(grpcServer, server)

	reflection.Register(grpcServer)

	serviceInfo := grpcServer.GetServiceInfo()
	for name, info := range serviceInfo {
		slog.Info("Service info", "name", name, "info", info)
	}
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			slog.Error("Failed to start grpc server", "err", err)
			return
		}
	}()
	return grpcServer
}

// Run starts gRPC server with provided config
func Run(cfg *Config) (*grpc.Server, error) {
	_, port, err := net.SplitHostPort(cfg.GrpcEndpoint)
	if err != nil {
		log.Fatalf("failed parse grpc endpoint, err: %v", err)
		return nil, err
	}
	// Create TCP listener
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return nil, err
	}
	log.Printf("Listening on %s\n", port)

	server, err := newServer(cfg)
	if err != nil {
		slog.Error("newServer failed", "error", err)
		return nil, err
	}

	// Wire both to gRPC
	grpcServer := setupRpcServer(lis, server)

	return grpcServer, nil
}

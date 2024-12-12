// Package resolver represents resolver node.
package resolver

import (
	"fmt"
	"log"
	"log/slog"
	"net"

	pb "github.com/1inch/p2p-network/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Run starts gRPC server with provided config
func Run(cfg *Config) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on %d\n", cfg.Port)
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	server, err := newServer(NewDefaultApiHandler())
	if err != nil {
		slog.Error("newServer failed", "error", err)
		return err
	}

	pb.RegisterExecuteServer(grpcServer, server)

	reflection.Register(grpcServer)

	serviceInfo := grpcServer.GetServiceInfo()
	for name, info := range serviceInfo {
		slog.Info("Service info", "name", name, "info", info)
	}
	err = grpcServer.Serve(lis)
	return err
}

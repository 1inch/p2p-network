// Package main represents resolver node.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/1inch/p2p-network/proto"
	"github.com/1inch/p2p-network/resolver/rpc"
	"google.golang.org/grpc"
)

func main() {
	var port = flag.Int("port", 8001, "Listener port")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on %d\n", *port)
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterExecuteServer(grpcServer, rpc.NewRpcServer())
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}

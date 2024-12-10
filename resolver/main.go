// Package main represents resolver node.
package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"

	"github.com/1inch/p2p-network/resolver/rpc"
)

func main() {
	var port = flag.Int("port", 8001, "Listener port")
	flag.Parse()
	lis, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", c.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on %d\n", c.Port)

	rpc.Start(&rpc.Config{Port: *port, LogLevel: slog.LevelInfo})
}

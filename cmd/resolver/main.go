// Package main implements cli wrapper for resolver
package main

import (
	"log/slog"
	"os"

	"github.com/1inch/p2p-network/resolver"
	"github.com/urfave/cli"
)

// TODO: setup cli interface
func main() {
	app := &cli.App{
		Name:  "resolver",
		Usage: "Resolver node",
		Commands: []cli.Command{
			{
				Name:  "run",
				Usage: "Runs resolver node",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "port",
						Value: 8001,
						Usage: "gRPC server port",
					},
				},
				Action: func(c *cli.Context) error {
					port := c.Int("port")
					err := resolver.Run(&resolver.Config{Port: port})
					if err != nil {
						slog.Error("Error starting server", "err", err)
						// TODO: handle error
					}

					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		slog.Error("Failed to run app", "err", err)
	}
}

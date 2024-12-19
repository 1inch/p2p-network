// Package main implements cli wrapper for resolver
package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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
					&cli.StringSliceFlag{
						Name:  "api",
						Value: &cli.StringSlice{"default"},
						Usage: "Supported APIs (default,infura)",
					},
					&cli.StringFlag{
						Name:   "infuraKey",
						Value:  "",
						Usage:  "Infura API Key",
						EnvVar: "INFURA_KEY",
					},
				},
				Action: func(c *cli.Context) error {
					port := c.Int("port")
					apis := c.StringSlice("api")
					var apiConfigs resolver.ApiConfigs
					for _, api := range apis {
						switch api {
						case "default":
							apiConfigs.Default.Enabled = true
						case "infura":
							apiConfigs.Infura.Enabled = true
							apiConfigs.Infura.Key = c.String("infuraKey")
						}
					}

					grpcServer, err := resolver.Run(&resolver.Config{Port: port, Apis: apiConfigs})
					if err != nil {
						slog.Error("Error starting server", "err", err)
						// TODO: handle error
						return err
					}
					interrupted := make(chan os.Signal, 1)
					signal.Notify(interrupted, syscall.SIGINT, syscall.SIGTERM)
					<-interrupted
					grpcServer.GracefulStop()

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

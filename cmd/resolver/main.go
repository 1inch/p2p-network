// Package main implements cli wrapper for resolver
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/1inch/p2p-network/internal/configs"
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
					&cli.StringFlag{
						Name:  "api",
						Value: "default",
						Usage: "Supported API (default,infura)",
					},
					&cli.StringFlag{
						Name:   "infuraKey",
						Value:  "",
						Usage:  "Infura API Key",
						EnvVar: "INFURA_KEY",
					},
					&cli.StringFlag{
						Name:  "configFile",
						Usage: "Path to the configuration file",
					},
				},
				Action: func(c *cli.Context) error {
					cfg := &resolver.Config{}
					// Try to load config file, if present
					configPath := c.String("configFile")
					if configPath != "" {
						cfgFromFile, err := configs.LoadConfig[resolver.Config](configPath)
						if err != nil {
							slog.Error("Error opening config file", "err", err)
						}
						cfg = cfgFromFile
					}
					// Override config file value for port
					port := c.Int("port")
					if port != 0 {
						cfg.Port = port
					}
					// Override config file value for apis
					api := c.String("api")
					if len(api) > 0 {
						var apiConfigs resolver.ApiConfigs
						switch api {
						case "default":
							apiConfigs.Default.Enabled = true
						case "infura":
							apiConfigs.Infura.Enabled = true
							apiConfigs.Infura.Key = c.String("infuraKey")
						}
						cfg.Apis = apiConfigs
					}
					grpcServer, err := resolver.Run(cfg)
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
			cliCommandRegistration(),
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		slog.Error("Failed to run app", "err", err)
	}
}

func cliCommandRegistration() cli.Command {
	return cli.Command{
		Name:  "registration",
		Usage: "Registration resolver in node registry",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "eth.addr",
				Value: "127.0.0.1",
				Usage: "address to ethereum node",
			},
			&cli.StringFlag{
				Name:  "eth.port",
				Value: "8545",
				Usage: "port to ethereum node",
			},
			&cli.StringFlag{
				Name:     "address",
				Usage:    "contract address where the register is located",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "account",
				Usage:    "account that will pay fee for adding resolver to registry",
				Required: true,
			},
			// TODO PrivateKey this volume cant be a command line parameter.
			// I think need move configurate private key from file in future.
			&cli.StringFlag{
				Name:     "privKey",
				Usage:    "account private key in hex which pay fee for register resolver",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "node.ip",
				Usage:    "this ip will set for resolver node",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "node.encodedPublicKey",
				Required: true,
				Usage:    "public key encoded in base64, will set for resolver node",
			},
		},
		Action: func(c *cli.Context) error {
			loggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			})
			logger := slog.New(loggerHandler)

			contractAddress := c.String("address")
			rawUrl := fmt.Sprintf("http://%s:%s", c.String("eth.addr"), c.String("eth.port"))
			logger.Info(rawUrl)
			regResolver, err := resolver.NewRegistrationResolver(logger, rawUrl, contractAddress)
			if err != nil {
				logger.Info("error when try create registration resolver", slog.Any("err", err.Error()))
				return err
			}

			ip := c.String("node.ip")
			account := c.String("account")
			privateKey := c.String("privKey")
			encodedPublicKey := c.String("node.encodedPublicKey")
			txHash, err := regResolver.Register(context.Background(), account, privateKey, ip, encodedPublicKey)
			if err != nil {
				return err
			}

			logger.Info("tx hash for registration new resolver", slog.Any("tx-hash", txHash))
			return nil
		},
	}
}

// Package main implements cli wrapper for resolver
package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/1inch/p2p-network/internal/configs"
	"github.com/1inch/p2p-network/resolver"
	"github.com/urfave/cli"
)

var (
	errContractAddressRequired = errors.New("contract address required")
	errRpcUrlRequired          = errors.New("rpc url required")
	errGrpcEndpointRequired    = errors.New("grpc endpoint required")
	errPrivateKeyRequired      = errors.New("private key required")
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
					&cli.StringFlag{
						Name:  "grpc_endpoint",
						Usage: "gRPC server endpoint",
					},
					&cli.StringFlag{
						Name:  "api",
						Value: "default",
						Usage: "Supported API (default,infura)",
					},
					&cli.StringFlag{
						Name:   "infura_key",
						Value:  "",
						Usage:  "Infura API Key",
						EnvVar: "INFURA_KEY",
					},
					&cli.StringFlag{
						Name:   "1inch_key",
						Value:  "",
						Usage:  "1Inch API key",
						EnvVar: "1INCH_KEY",
					},
					&cli.StringFlag{
						Name:  "config_file",
						Usage: "Path to the configuration file",
					},
				},
				Action: func(c *cli.Context) error {
					cfg := resolver.Config{}
					// Try to load config file, if present
					loadedCfg := loadConfigByPath(c.String("config_file"))
					if loadedCfg != nil {
						cfg = *loadedCfg
					}
					logger := setupLogger(cfg)

					// Override config file value for port
					grpcEndpoint := c.String("grpc_endpoint")
					if grpcEndpoint != "" {
						cfg.GrpcEndpoint = grpcEndpoint
					}
					// Override config file value for apis
					if !isApiHandlerSet(&cfg) {
						api := c.String("api")
						if len(api) > 0 {
							var apiConfigs resolver.ApiConfigs
							switch api {
							case "default":
								apiConfigs.Default.Enabled = true
							case "infura":
								apiConfigs.Infura.Enabled = true
								apiConfigs.Infura.Key = c.String("infura_key")
							case "1inch":
								apiConfigs.OneInch.Enabled = true
								apiConfigs.OneInch.Key = c.String("1inch_key")
							}
							cfg.Apis = apiConfigs
						}
					}

					resolverNode, err := resolver.New(cfg, logger)
					if err != nil {
						logger.Error("Error create resolver node", slog.Any("err", err))
						return err
					}
					err = resolverNode.Run()
					if err != nil {
						logger.Error("Failed start resolver", slog.Any("err", err))
						return err
					}

					interrupted := make(chan os.Signal, 1)
					signal.Notify(interrupted, syscall.SIGINT, syscall.SIGTERM)
					<-interrupted
					return resolverNode.Stop()
				},
			},
			cliCommandRegister(),
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		slog.Error("Failed to run app", "err", err)
	}
}

func cliCommandRegister() cli.Command {
	return cli.Command{
		Name:  "register",
		Usage: "Register resolver in node registry",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "rpc_url",
				Usage: "rpc url to blockchain node",
			},
			&cli.StringFlag{
				Name:  "contract_address",
				Usage: "contract address where the register is located",
			},
			&cli.StringFlag{
				Name:  "private_key",
				Usage: "account private key in hex which pay fee for register resolver",
			},
			&cli.StringFlag{
				Name:  "grpc_endpoint",
				Usage: "this endpoint will set for resolver node",
			},
			&cli.StringFlag{
				Name:  "config_file",
				Usage: "Path to the configuration file",
			},
		},
		Action: func(c *cli.Context) error {
			loggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			})
			logger := slog.New(loggerHandler)
			cfg := &resolver.Config{}

			// Try to load config file, if present
			loadedCfg := loadConfigByPath(c.String("config_file"))
			if loadedCfg != nil {
				cfg = loadedCfg
			} else {
				// check if this param is not null and override
				contractAddr := c.String("contract_address")
				if contractAddr != "" {
					cfg.ContractAddress = contractAddr
				} else {
					return errContractAddressRequired
				}

				// check if this param is not null and override
				rpc_url := c.String("rpc_url")
				if rpc_url != "" {
					cfg.RpcUrl = rpc_url
				} else {
					return errRpcUrlRequired
				}

				grpc_endpoint := c.String("grpc_endpoint")
				if grpc_endpoint != "" {
					cfg.GrpcEndpoint = grpc_endpoint
				} else {
					return errGrpcEndpointRequired
				}

				privKey := c.String("private_key")
				if privKey != "" {
					cfg.PrivateKey = privKey
				} else {
					return errPrivateKeyRequired
				}
			}

			regResolver, err := resolver.NewRegistrationResolver(logger, cfg)
			if err != nil {
				logger.Info("error when try create registration resolver", slog.Any("err", err.Error()))
				return err
			}

			txHash, err := regResolver.Register(context.Background())
			if err != nil {
				return err
			}

			logger.Info("tx hash for registration new resolver", slog.Any("tx-hash", txHash))
			return nil
		},
	}
}

func loadConfigByPath(configPath string) *resolver.Config {
	if configPath != "" {
		cfgFromFile, err := configs.LoadConfig[resolver.Config](configPath)
		if err != nil {
			slog.Error("Error opening config file", "err", err)
			return nil
		}
		return cfgFromFile
	}

	return nil
}

func setupLogger(cfg resolver.Config) *slog.Logger {
	leveler := new(slog.LevelVar)
	leveler.Set(cfg.LogLevel.Level())
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: leveler,
	})
	return slog.New(handler)
}

func isApiHandlerSet(cfg *resolver.Config) bool {
	return cfg.Apis.Default.Enabled || cfg.Apis.Infura.Enabled || cfg.Apis.OneInch.Enabled
}

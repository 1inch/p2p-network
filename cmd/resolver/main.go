// Package main implements cli wrapper for resolver
package main

import (
	"context"
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
					loadedCfg := loadConfigByPath(c.String("configFile"))
					if loadedCfg != nil {
						cfg = loadedCfg
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
					grpcServer, _, err := resolver.Run(cfg)
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
				Name:  "privKey",
				Usage: "account private key in hex which pay fee for register resolver",
			},
			&cli.StringFlag{
				Name:  "ip",
				Usage: "this ip will set for resolver node",
			},
			&cli.StringFlag{
				Name:  "configFile",
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
			loadedCfg := loadConfigByPath(c.String("configFile"))
			if loadedCfg != nil {
				cfg = loadedCfg
			}

			// check if this param is not null and override
			contractAddr := c.String("contract_address")
			if contractAddr != "" {
				cfg.ContractAddress = contractAddr
			}

			// check if this param is not null and override
			rpc_url := c.String("rpc_url")
			if rpc_url != "" {
				cfg.RpcUrl = rpc_url
			}

			ip := c.String("ip")
			if ip != "" {
				cfg.Ip = ip
			}

			privKey := c.String("privKey")
			if privKey != "" {
				cfg.PrivateKey = privKey
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

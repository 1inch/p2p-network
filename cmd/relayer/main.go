// Package main contains the entrypoint for the relayer node.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/1inch/p2p-network/internal/configs"
	"github.com/1inch/p2p-network/internal/log"
	"github.com/1inch/p2p-network/relayer"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:  "relayer",
		Usage: "Relayer node",
		Commands: []cli.Command{
			{
				Name:  "run",
				Usage: "Runs the relayer node",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Usage:    "Path to the configuration file",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					leveler := new(slog.LevelVar)
					leveler.Set(slog.LevelInfo)
					handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
						Level: leveler,
					})
					logger := slog.New(handler)
					logger.Info("starting relayer node")

					configPath := c.String("config")
					cfg, err := configs.LoadConfig[relayer.Config](configPath)
					if err != nil {
						logger.Error("failed to load relayer node configuration", slog.String("path", configPath), slog.Any("err", err))
					}
					logLevel, err := log.ParseLevel(cfg.LogLevel)
					if err != nil {
						logger.Error("failed parsing log level", slog.String("log_level", cfg.LogLevel), slog.Any("err", err))
					}
					leveler.Set(logLevel)

					logger.Info("config file loaded", slog.String("path", configPath))

					node, err := relayer.New(cfg, logger)
					if err != nil {
						logger.Error("failed to initialize relayer node", slog.Any("err", err))
					}

					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()

					go handleInterrupt(cancel)

					if err := node.Run(ctx); err != nil {
						logger.Error("failed to run relayer node", slog.Any("err", err))
					}

					logger.Info("relayer node stopped gracefully")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		handler := slog.NewTextHandler(os.Stdout, nil)
		logger := slog.New(handler)
		logger.Error("failed to run relayer node CLI interface", slog.Any("err", err))
	}
}

func handleInterrupt(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	cancel()
}

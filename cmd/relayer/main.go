// Package main contains the entrypoint for the relayer node.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/1inch/p2p-network/internal/configs"
	"github.com/1inch/p2p-network/relayer"
	"github.com/sirupsen/logrus"
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
					logger := logrus.New()
					logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
					logger.SetLevel(logrus.DebugLevel)
					logger.Info("starting relayer node")

					configPath := c.String("config")
					cfg, err := configs.LoadConfig[relayer.Config](configPath)
					if err != nil {
						logger.WithError(err).WithField("path", configPath).Error("failed to load relayer node configuration")
					}

					logger.WithField("path", configPath).Info("config file loaded")

					node, err := relayer.New(cfg, logger)
					if err != nil {
						logger.WithError(err).Error("failed to initialize relayer node")
					}

					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()

					// TODO: handle interrupts?
					go handleInterrupt(cancel)

					if err := node.Run(ctx); err != nil {
						logger.WithError(err).Error("failed to run relayer node")
					}

					logger.Info("relayer node stopped gracefully")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger := logrus.New()
		logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
		logger.WithError(err).Error("failed to run relayer node CLI interface")
	}
}

func handleInterrupt(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	cancel()
}

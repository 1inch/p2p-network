// Package relayer represents relayer node.
package relayer

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"sync"

	"github.com/1inch/p2p-network/relayer/httpapi"
	"github.com/1inch/p2p-network/relayer/webrtc"
	"golang.org/x/sync/errgroup"
)

// Relayer represents the core relayer node with subsystems.
type Relayer struct {
	Config       *Config
	Logger       *slog.Logger
	WebRTCServer *webrtc.Server
	HTTPServer   *httpapi.Server
	wg           sync.WaitGroup
}

// New initializes a new Relayer instance with provided configuration and logger.
func New(cfg *Config, logger *slog.Logger) (*Relayer, error) {
	var httpServer *httpapi.Server
	{
		// setup http listener.
		httpListener, err := net.Listen("tcp4", cfg.HTTPEndpoint)
		if err != nil {
			logger.Error("http server failed to listen on tcp", slog.String("addr", cfg.HTTPEndpoint), slog.Any("err", err))
			return nil, err
		}
		mux := http.NewServeMux()
		// TODO: add handlers
		httpServer = httpapi.New(logger.WithGroup("httpapi"), httpListener, mux)
	}

	return &Relayer{
		Config:     cfg,
		Logger:     logger,
		HTTPServer: httpServer,
	}, nil
}

// Run starts the relayer and its subsystems.
func (r *Relayer) Run(ctx context.Context) error {
	var group errgroup.Group
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	group.Go(func() error {
		defer cancel()

		r.Logger.Info("http server started", slog.String("addr", r.HTTPServer.Addr()))
		err := r.HTTPServer.Run(childCtx)
		if err != nil {
			r.Logger.Error("http server failed to serve", slog.Any("err", err))
			return err
		}

		return nil
	})

	// Wait for all goroutines to complete or an error to occur
	if err := group.Wait(); err != nil {
		r.Logger.Error("relayer encountered an error", slog.Any("err", err))
		return err
	}

	return nil
}

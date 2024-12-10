// Package relayer represents relayer node.
package relayer

import (
	"context"
	"net"
	"net/http"
	"sync"

	"github.com/1inch/p2p-network/relayer/httpapi"
	"github.com/1inch/p2p-network/relayer/webrtc"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// Relayer represents the core relayer node with subsystems.
type Relayer struct {
	Config       *Config
	Logger       *logrus.Logger
	WebRTCServer *webrtc.Server
	HTTPServer   *httpapi.Server
	wg           sync.WaitGroup
}

// New initializes a new Relayer instance with provided configuration and logger.
func New(cfg *Config, logger *logrus.Logger) (*Relayer, error) {
	var httpServer *httpapi.Server
	{
		// setup http listener.
		httpListener, err := net.Listen("tcp4", cfg.HTTPEndpoint)
		if err != nil {
			logger.WithError(err).WithField("addr", cfg.HTTPEndpoint).Error("http server failed to listen on tcp")
			return nil, err
		}
		mux := http.NewServeMux()
		// TODO: add handlers
		httpServer = httpapi.New(logger, httpListener, mux)
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

		r.Logger.WithField("addr", r.HTTPServer.Addr()).Info("http server started")
		err := r.HTTPServer.Run(childCtx)
		if err != nil {
			r.Logger.WithError(err).Error("http server failed to serve")
			return err
		}

		return nil
	})

	// Wait for all goroutines to complete or an error to occur
	if err := group.Wait(); err != nil {
		r.Logger.WithError(err).Error("relayer encountered an error")
		return err
	}

	return nil
}

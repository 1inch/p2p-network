// Package relayer represents relayer node.
package relayer

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"sync"

	"github.com/1inch/p2p-network/relayer/grpc"
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
	sdpRequests := make(chan webrtc.SDPRequest)
	var httpServer *httpapi.Server
	{
		// setup http listener.
		httpListener, err := net.Listen("tcp4", cfg.HTTPEndpoint)
		if err != nil {
			logger.Error("http server failed to listen on tcp", slog.String("addr", cfg.HTTPEndpoint), slog.Any("err", err))
			return nil, err
		}
		mux := http.NewServeMux()
		mux.HandleFunc("POST /sdp", webrtc.SDPHandler(logger, sdpRequests))
		httpServer = httpapi.New(logger.WithGroup("httpapi"), httpListener, mux)
	}

	var werbrtcServer *webrtc.Server
	{
		// setup webrtc listener.
		var err error
		grpcClient, err := grpc.New(cfg.GRPCServerAddress)
		if err != nil {
			logger.Error("failed to initialize grpc client", slog.Any("err", err))
			return nil, err
		}
		werbrtcServer, err = webrtc.New(logger.WithGroup("webrtc"), cfg.WebRTCICEServer, grpcClient, sdpRequests)
		if err != nil {
			logger.Error("failed to create webrtc server", slog.String("iceserver", cfg.WebRTCICEServer), slog.Any("err", err))
			return nil, err
		}
	}

	return &Relayer{
		Config:       cfg,
		Logger:       logger,
		HTTPServer:   httpServer,
		WebRTCServer: werbrtcServer,
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

	group.Go(func() error {
		defer cancel()

		r.Logger.Info("webrtc server started", slog.String("iceserver", r.Config.WebRTCICEServer))
		err := r.WebRTCServer.Run(childCtx)
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

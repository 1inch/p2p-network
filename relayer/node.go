// Package relayer represents relayer node.
package relayer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/1inch/p2p-network/contracts"
	"github.com/1inch/p2p-network/relayer/grpc"
	"github.com/1inch/p2p-network/relayer/httpapi"
	"github.com/1inch/p2p-network/relayer/webrtc"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"golang.org/x/sync/errgroup"
)

// Relayer represents the core relayer node with subsystems.
type Relayer struct {
	Config       *Config
	Logger       *slog.Logger
	WebRTCServer *webrtc.Server
	HTTPServer   *httpapi.Server
}

// New initializes a new Relayer instance with provided configuration and logger.
func New(cfg *Config, logger *slog.Logger) (*Relayer, error) {
	sdpRequests := make(chan webrtc.SDPRequest)
	iceCandidates := make(chan webrtc.ICECandidate)
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
		mux.HandleFunc("POST /candidate", webrtc.CandidateHandler(logger, iceCandidates))
		mux.HandleFunc("GET /relayer", func(w http.ResponseWriter, r *http.Request) {
			client, err := contracts.Dial(r.Context(), cfg.BlockchainRPCAddress, cfg.PrivateKey, cfg.ContractAddress)
			if err != nil {
				http.Error(w, "failed to connect to Ethereum node", http.StatusInternalServerError)
				return
			}
			defer client.Close()

			ip, err := client.Registry.GetRelayer(&bind.CallOpts{})
			if err != nil {
				http.Error(w, "failed to get closest relayer node", http.StatusInternalServerError)
				return
			}

			resp := struct {
				IPAddress string `json:"ip_address"`
			}{IPAddress: ip}

			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(resp)
			if err != nil {
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
				return
			}
		})
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
		werbrtcServer, err = webrtc.New(logger.WithGroup("webrtc"), cfg.WebRTCICEServer, grpcClient, sdpRequests, iceCandidates)
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
		r.Logger.Info("register relayer with node registry", slog.String("ip_address", r.HTTPServer.Addr()))
		err := r.registerRelayer(childCtx)
		if err != nil {
			r.Logger.Error("failed to register relayer with node registry", slog.Any("err", err))
			return err
		}

		return nil
	})

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
		r.Logger.Info("webrtc server started", slog.String("iceserver", r.Config.WebRTCICEServer))
		err := r.WebRTCServer.Run(childCtx)
		if err != nil {
			r.Logger.Error("webrtc server failed to serve", slog.Any("err", err))
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

func (r *Relayer) registerRelayer(ctx context.Context) error {
	client, err := contracts.Dial(ctx, r.Config.BlockchainRPCAddress, r.Config.PrivateKey, r.Config.ContractAddress)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}
	defer client.Close()

	tx, err := client.Registry.RegisterRelayer(client.Auth, r.Config.HTTPEndpoint)
	if err != nil {
		return fmt.Errorf("failed to register relayer: %w", err)
	}

	if err := client.WaitForTx(ctx, tx.Hash()); err != nil {
		return fmt.Errorf("wait for transaction error: %w", err)
	}

	r.Logger.Debug("successfully sent register relayer tx", slog.String("tx_hash", tx.Hash().String()))

	return nil
}

// Package relayer represents relayer node.
package relayer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/1inch/p2p-network/internal/registry"
	"github.com/1inch/p2p-network/relayer/grpc"
	"github.com/1inch/p2p-network/relayer/httpapi"
	webrtcserver "github.com/1inch/p2p-network/relayer/webrtc"
	"github.com/pion/webrtc/v4"
	"golang.org/x/sync/errgroup"
)

// Relayer represents the core relayer node with subsystems.
type Relayer struct {
	Config       *Config
	Logger       *slog.Logger
	WebRTCServer *webrtcserver.Server
	HTTPServer   *httpapi.Server
}

// New initializes a new Relayer instance with provided configuration and logger.
func New(cfg *Config, logger *slog.Logger) (*Relayer, error) {
	if len(cfg.WebrtcConfig.ICEServers) == 0 {
		return nil, webrtcserver.ErrInvalidICEServer
	}

	sdpRequests := make(chan webrtcserver.SDPRequest)
	iceCandidates := make(chan webrtcserver.ICECandidate)
	var httpServer *httpapi.Server
	{
		// setup http listener.
		httpListener, err := net.Listen("tcp4", cfg.HTTPEndpoint)
		if err != nil {
			logger.Error("http server failed to listen on tcp", slog.String("addr", cfg.HTTPEndpoint), slog.Any("err", err))
			return nil, err
		}
		mux := http.NewServeMux()
		mux.HandleFunc("POST /sdp", webrtcserver.SDPHandler(logger, sdpRequests))
		mux.HandleFunc("POST /candidate", webrtcserver.CandidateHandler(logger, iceCandidates))
		mux.HandleFunc("GET /relayer", func(w http.ResponseWriter, r *http.Request) {
			client, err := registry.Dial(r.Context(), &registry.Config{
				DialURI:         cfg.DiscoveryConfig.RpcUrl,
				PrivateKey:      cfg.PrivateKey,
				ContractAddress: cfg.DiscoveryConfig.ContractAddress,
			})
			if err != nil {
				http.Error(w, "failed to connect to Ethereum node", http.StatusInternalServerError)
				return
			}
			defer client.Close()

			ip, resolvers, err := client.GetRelayer()
			if err != nil {
				http.Error(w, "failed to get closest relayer node", http.StatusInternalServerError)
				return
			}

			resp := struct {
				IPAddress string   `json:"ip_address"`
				Resolvers [][]byte `json:"resolvers"`
			}{IPAddress: ip, Resolvers: resolvers}

			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(resp)
			if err != nil {
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
				return
			}
		})
		mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
			logger.Debug("called /health endpoint")
		})

		httpServer = httpapi.New(logger.WithGroup("httpapi"), httpListener, corsMiddleware(mux))
	}

	var werbrtcServer *webrtcserver.Server
	{
		// setup webrtc listener.
		var err error
		ctx := context.Background()
		registryClient, err := registry.Dial(ctx, &registry.Config{
			DialURI:         cfg.DiscoveryConfig.RpcUrl,
			PrivateKey:      cfg.PrivateKey,
			ContractAddress: cfg.DiscoveryConfig.ContractAddress,
		})
		if err != nil {
			logger.Error("failed to initialize registry client", slog.Any("err", err))
			return nil, err
		}
		werbrtcServer, err = webrtcserver.New(logger.WithGroup("webrtc"), iceServerByConfig(*cfg), grpc.New(registryClient), sdpRequests, iceCandidates, webrtcOptionsByConfig(*cfg)...)
		if err != nil {
			logger.Error("failed to create webrtc server", slog.Any("err", err))
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
		r.Logger.Info("webrtc server started")
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

// RegisterRelayer registers the relayer node with the registry contract.
func (r *Relayer) RegisterRelayer(ctx context.Context) error {
	client, err := registry.Dial(ctx, &registry.Config{
		DialURI:         r.Config.DiscoveryConfig.RpcUrl,
		PrivateKey:      r.Config.PrivateKey,
		ContractAddress: r.Config.DiscoveryConfig.ContractAddress,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}
	defer client.Close()

	if err = client.RegisterRelayer(ctx, r.Config.HTTPEndpoint); err != nil {
		return fmt.Errorf("failed to register relayer: %w", err)
	}

	return nil
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func iceServerByConfig(cfg Config) []webrtc.ICEServer {
	iceServers := make([]webrtc.ICEServer, len(cfg.WebrtcConfig.ICEServers))

	for i := range cfg.WebrtcConfig.ICEServers {
		iceServers[i] = webrtc.ICEServer{
			URLs:           []string{cfg.WebrtcConfig.ICEServers[i].Url},
			Username:       cfg.WebrtcConfig.ICEServers[i].Username,
			Credential:     cfg.WebrtcConfig.ICEServers[i].Password,
			CredentialType: webrtc.ICECredentialTypePassword,
		}
	}

	return iceServers
}

func webrtcOptionsByConfig(cfg Config) []webrtcserver.Option {
	opts := []webrtcserver.Option{}

	if cfg.WebrtcConfig.UseTrickleICE {
		opts = append(opts, webrtcserver.WithTrickleICE())
	}
	if cfg.WebrtcConfig.RetryConfig.Enabled {
		opts = append(opts, webrtcserver.WithRetry(webrtcserver.Retry{
			Count:    cfg.WebrtcConfig.RetryConfig.Count,
			Interval: cfg.WebrtcConfig.RetryConfig.Interval,
		}))
	}
	if cfg.WebrtcConfig.PeerPortConfig.Enabled {
		opts = append(opts, webrtcserver.WithPeerPort(webrtcserver.PeerRangePort{
			Min: cfg.WebrtcConfig.PeerPortConfig.Min,
			Max: cfg.WebrtcConfig.PeerPortConfig.Max,
		}))
	}

	return opts
}

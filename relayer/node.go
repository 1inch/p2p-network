// Package relayer represents relayer node.
package relayer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"

	"github.com/1inch/p2p-network/internal/registry"
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
			client, err := registry.Dial(r.Context(), &registry.Config{
				DialURI:         cfg.BlockchainRPCAddress,
				PrivateKey:      cfg.PrivateKey,
				ContractAddress: cfg.ContractAddress,
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

		httpServer = httpapi.New(logger.WithGroup("httpapi"), httpListener, handlerWithLoggingAndCors(logger, mux))
	}

	var werbrtcServer *webrtc.Server
	{
		// setup webrtc listener.
		var err error
		ctx := context.Background()
		registryClient, err := registry.Dial(ctx, &registry.Config{
			DialURI:         cfg.BlockchainRPCAddress,
			PrivateKey:      cfg.PrivateKey,
			ContractAddress: cfg.ContractAddress,
		})
		if err != nil {
			logger.Error("failed to initialize registry client", slog.Any("err", err))
			return nil, err
		}

    webrtcCfg := webrtc.RetryRequestConfig{
			Count:    cfg.RetryRequestConfig.Count,
			Interval: cfg.RetryRequestConfig.Interval,
		}
		werbrtcServer, err = webrtc.New(webrtcCfg, logger.WithGroup("webrtc"), cfg.WebRTCICEServer, grpc.New(registryClient), sdpRequests, iceCandidates)

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

// RegisterRelayer registers the relayer node with the registry contract.
func (r *Relayer) RegisterRelayer(ctx context.Context) error {
	client, err := registry.Dial(ctx, &registry.Config{
		DialURI:         r.Config.BlockchainRPCAddress,
		PrivateKey:      r.Config.PrivateKey,
		ContractAddress: r.Config.ContractAddress,
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

func corsMiddleware(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
	}
}

func loggerMiddlewareBeforeServe(logger *slog.Logger, r *http.Request) *http.Request {
	// cloned http request because only once can be read request body
	clonedReq := r.Clone(r.Context())
	logger.Info("receive request", slog.Any("method", r.Method), slog.Any("api", r.URL.Path))
	logger.Debug("with request", slog.Any("body", func() string {
		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Debug("failed get request bytes from reader")
				return ""
			}

			clonedReq.Body = io.NopCloser(bytes.NewReader(bodyBytes))

			return string(bodyBytes)
		}
		return "<nil>"
	}()))

	return clonedReq
}

func loggerMiddlewareAfterServe(logger *slog.Logger, r *http.Request) {
	logger.Info("request process", slog.Any("method", r.Method), slog.Any("api", r.URL.Path))
}

func handlerWithLoggingAndCors(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corsMiddleware(w, r)
		clonedReq := loggerMiddlewareBeforeServe(logger, r)

		next.ServeHTTP(w, clonedReq)
		loggerMiddlewareAfterServe(logger, clonedReq)
	})
}

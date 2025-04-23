// Package resolver represents resolver node.
package resolver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"bursavich.dev/grpcprom"
	pb "github.com/1inch/p2p-network/proto/resolver"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var errNoHandlerApiInConfig = errors.New("no handler api in config")

// Resolver represents node with subsystems.
type Resolver struct {
	cfg              Config
	logger           *slog.Logger
	httpMetricServer *http.Server
	grpcServer       *grpc.Server
	lis              net.Listener
}

// New method for create new instance of Resolver
func New(cfg Config, logger *slog.Logger) (*Resolver, error) {
	resolver := &Resolver{
		cfg:    cfg,
		logger: logger,
	}

	server, err := newServer(&cfg)
	if err != nil {
		logger.Error("failed create server", slog.Any("err", err.Error()))
	}

	loggingInterceptor := func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		return loggingRequestHandler(ctx, logger, req, info, handler)
	}
	var serverOpts []grpc.ServerOption
	// if metric enabled setup http server for metrics
	registry := prometheus.NewRegistry()
	serverMetrics := grpcprom.NewServerMetrics()

	if cfg.Metric.Enabled {
		serverMetrics := grpcprom.NewServerMetrics()
		registry.MustRegister(serverMetrics)

		if err != nil {
			logger.Error("failed to start prometheus exporter", slog.Any("err", err.Error()))
			return nil, err
		}

		serverOpts = append(serverOpts,
			grpc.ChainUnaryInterceptor(loggingInterceptor, serverMetrics.UnaryInterceptor()),
			grpc.StatsHandler(serverMetrics.StatsHandler()),
		)

		metricServer := newMetricServer(&cfg, registry)
		resolver.httpMetricServer = metricServer
	} else {
		serverOpts = append(serverOpts, grpc.UnaryInterceptor(loggingInterceptor))
	}

	resolver.grpcServer = newGrpcServer(logger, server, serverOpts...)
	serverMetrics.Init(resolver.grpcServer)

	// Create TCP listener
	listener, err := net.Listen("tcp", cfg.GrpcEndpoint)
	if err != nil {
		logger.Error("Failed to create net listener", slog.Any("err", err.Error()))
		return nil, err
	}
	resolver.lis = listener

	return resolver, nil
}

func newGrpcServer(logger *slog.Logger, server *Server, opts ...grpc.ServerOption) *grpc.Server {
	grpcServer := grpc.NewServer(opts...)
	healthServer := health.NewServer()

	pb.RegisterExecuteServer(grpcServer, server)
	grpchealth.RegisterHealthServer(grpcServer, healthServer)

	// TODO maybe need make this turn on/off by configuration?
	reflection.Register(grpcServer)

	serviceInfo := grpcServer.GetServiceInfo()
	for name := range serviceInfo {
		logger.Info("service info", slog.Any("name", name))
	}

	return grpcServer
}

func newMetricServer(cfg *Config, registry prometheus.Gatherer) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("GET /metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	metricServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Metric.Port),
		Handler: mux,
	}

	return metricServer
}

// Run starts gRPC server with provided config
func (r *Resolver) Run() error {
	go func() {
		r.logger.Info("listening grpc server", slog.Any("address", r.Addr()))
		if err := r.grpcServer.Serve(r.lis); err != nil {
			r.logger.Error("failed to start grpc server", slog.Any("err", err.Error()))
			return
		}
	}()

	if r.cfg.Metric.Enabled {
		go func() {
			r.logger.Info("listening metric server", slog.Any("port", r.cfg.Metric.Port))
			if err := r.httpMetricServer.ListenAndServe(); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					r.logger.Info("metric server closed")
				} else {
					r.logger.Error("failed to start http server", slog.Any("err", err.Error()))
				}
			}
		}()
	}

	return nil
}

// Stop represented method for graceful stop Resolver and
func (r *Resolver) Stop() error {
	r.grpcServer.GracefulStop()

	if r.cfg.Metric.Enabled {
		ctx := context.Background()
		err := r.httpMetricServer.Shutdown(ctx)
		if err != nil {
			r.logger.Error("failed shutdown http metric server")
			return err
		}
	}

	r.logger.Info("grpc server stopped")
	return nil
}

func loggingRequestHandler(ctx context.Context, logger *slog.Logger, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logger.Info("received request on grpc server", slog.Any("method", info.FullMethod))
	logger.Debug("with request", slog.Any("req", protojson.Format(req.(proto.Message))))

	resp, err := handler(ctx, req)

	if err != nil {
		logger.Info("request failed process", slog.Any("method", info.FullMethod))
		logger.Debug("with error", slog.Any("err", err.Error()))
	} else {
		logger.Info("request process success", slog.Any("method", info.FullMethod))
		logger.Debug("with response", slog.Any("resp", protojson.Format(resp.(proto.Message))))
	}

	return resp, err
}

// Addr returns the net listener address.
func (r *Resolver) Addr() string {
	return r.lis.Addr().String()
}

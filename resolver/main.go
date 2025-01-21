// Package resolver represents resolver node.
package resolver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	pb "github.com/1inch/p2p-network/proto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/stats/opentelemetry"
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
		logger.Error("Failed create server", slog.Any("err", err.Error()))
	}

	// if metric enabled setup http server for metrics
	if cfg.Metric.Enabled {
		exporter, err := prometheus.New()
		if err != nil {
			logger.Error("Failed to start prometheus exporter", slog.Any("err", err.Error()))
			return nil, err
		}

		meterProvider := metric.NewMeterProvider(metric.WithReader(exporter))
		meterServerOption := opentelemetry.ServerOption(opentelemetry.Options{
			MetricsOptions: opentelemetry.MetricsOptions{
				MeterProvider: meterProvider,
			},
		})
		serverOption := meterServerOption
		metricServer := newMetricServer(&cfg)

		resolver.httpMetricServer = metricServer
		resolver.grpcServer = newGrpcServer(logger, server, serverOption)
	}

	resolver.grpcServer = newGrpcServer(logger, server)

	// Create TCP listener
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.Port))
	if err != nil {
		logger.Error("Failed to create net listener", slog.Any("err", err.Error()))
		return nil, err
	}
	resolver.lis = listener

	return resolver, nil
}

func newGrpcServer(logger *slog.Logger, server *Server, opts ...grpc.ServerOption) *grpc.Server {
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterExecuteServer(grpcServer, server)
	reflection.Register(grpcServer)

	serviceInfo := grpcServer.GetServiceInfo()
	for name := range serviceInfo {
		logger.Info("Service info", slog.Any("name", name))
	}

	return grpcServer
}

func newMetricServer(cfg *Config) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("GET /metrics", promhttp.Handler())

	metricServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Metric.Port),
		Handler: mux,
	}

	return metricServer
}

// Run starts gRPC server with provided config
func (r *Resolver) Run() error {
	go func() {
		r.logger.Info("Listening grpc server", slog.Any("port", r.cfg.Port))
		if err := r.grpcServer.Serve(r.lis); err != nil {
			r.logger.Error("Failed to start grpc server", slog.Any("err", err.Error()))
			return
		}
	}()

	if r.cfg.Metric.Enabled {
		go func() {
			r.logger.Info("Listening metric server", slog.Any("port", r.cfg.Metric.Port))
			if err := r.httpMetricServer.ListenAndServe(); err != nil {
				r.logger.Error("Failed to start http server", slog.Any("err", err.Error()))
				return
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
			r.logger.Error("Failed shutdown http metric server")
			return err
		}
	}
	return nil
}

// Addr returns the net listener address.
func (r *Resolver) Addr() string {
	return r.lis.Addr().String()
}

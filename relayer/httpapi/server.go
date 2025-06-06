// Package httpapi starts http server.
package httpapi

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

// Server wraps the http.Server and logger.
type Server struct {
	logger *slog.Logger
	lis    net.Listener
	server *http.Server
}

// New creates a new Server instance.
func New(logger *slog.Logger, lis net.Listener, handler http.Handler) *Server {
	server := &http.Server{
		Handler: handler,
	}

	return &Server{
		logger: logger,
		lis:    lis,
		server: server,
	}
}

// Run starts the HTTP server and serves requests.
func (s *Server) Run(ctx context.Context) error {
	var g errgroup.Group

	// Start serving HTTP requests
	g.Go(func() error {
		if err := s.server.Serve(s.lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("http server error", slog.String("addr", s.Addr()), slog.Any("err", err))
			return err
		}
		return nil
	})

	// Listen for shutdown signal
	g.Go(func() error {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		s.logger.Info("shutting down http server")

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("error shutting down http server", slog.Any("err", err))
			return err
		}

		s.logger.Info("http server shutdown completed")
		return nil
	})

	return g.Wait()
}

// Addr returns the HTTP server address.
func (s *Server) Addr() string {
	return s.lis.Addr().String()
}

package http

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

type Server struct {
	*http.Server
	stopped chan struct{}
}

func New(ctx context.Context, addr string, handler http.Handler) *Server {
	s := &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
	return s
}

// Wait must be called if Run have been called to re-sync.
func (s *Server) Wait() {
	<-s.stopped
	s.stopped = nil
}

var ErrServerRunAlreadyRunning = errors.New("HTTP server already launched")

// Run require calling Wait to ensure asynchronous internal processing finishes.
// It is mandatory to call Wait between each call to Run to restart.
func (s *Server) Run(ctx context.Context) error {
	if s.stopped != nil {
		return ErrServerRunAlreadyRunning
	}
	serverExited := make(chan struct{})
	go func() {
		defer close(serverExited)
		switch err := s.ListenAndServe(); err {
		case http.ErrServerClosed:
			zap.L().Info("HTTP server closed",
				zap.String("server_addr", s.Addr))
		case nil:
		default:
			zap.L().Error("HTTP server exit unexpected",
				zap.Error(err),
				zap.String("server_addr", s.Addr))
		}
	}()
	done := make(chan struct{})
	go func() {
		defer close(done)
		select {
		case <-ctx.Done():
			zap.L().Info("closing http server",
				zap.String("server_addr", s.Addr))
			timeout := config.GetDuration(serverShutdownTimeoutConfig)
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			switch err := s.Shutdown(ctx); err {
			case nil:
			default:
				zap.L().Info("HTTP server graceful shutdown",
					zap.Error(err),
					zap.String("server_addr", s.Addr),
					zap.Duration("exit_timeout", timeout))
				if err := s.Server.Close(); err != nil {
					zap.L().Info("HTTP server force exit",
						zap.Error(err),
						zap.String("server_addr", s.Addr))
				}
			}
			<-serverExited
		case <-serverExited:
		}
	}()
	zap.L().Info("HTTP server exited",
		zap.String("server_addr", s.Addr))
	s.stopped = done
	return nil
}

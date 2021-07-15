package http

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

type Server struct {
	*http.Server
	stop func()
}

func New(ctx context.Context, addr string, handler http.Handler) *Server {
	s := &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
	ctx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	go func() {
		defer close(done)
		s.run(ctx)
	}()
	s.stop = func() { cancel(); <-done }
	return s
}

func (s *Server) run(ctx context.Context) {
	serverExited := make(chan struct{})
	go func() {
		defer close(serverExited)
		switch err := s.ListenAndServe(); err {
		case http.ErrServerClosed:
			zap.L().Info("http server closed",
				zap.String("server_addr", s.Addr))
		case nil:
		default:
			zap.L().Error("http server exit unexpected",
				zap.Error(err),
				zap.String("server_addr", s.Addr))
		}
	}()
	select {
	case <-ctx.Done():
		zap.L().Info("closing http server",
			zap.String("server_addr", s.Addr))
		s.Close()
		<-serverExited
	case <-serverExited:
	}
	zap.L().Info("http server exited",
		zap.String("server_addr", s.Addr))
}

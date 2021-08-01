package main

import (
	"net"

	"github.com/7d4b9/zouti/context"
	"github.com/7d4b9/zouti/http"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var config = viper.New()

func main() {
	addr := net.JoinHostPort("", config.GetString("http_port"))
	ctx := context.CancelOnSigInterrupt
	s := http.New(ctx, addr, HTTPMux())
	if err := s.Run(ctx); err != nil {
		zap.L().Fatal("http server start",
			zap.String("addr", addr))
	}
	defer s.Wait()
}

package main

import (
	"net"

	"github.com/7d4b9/lever/context"
	"github.com/7d4b9/lever/http"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var config = viper.New()

func main() {
	addr := net.JoinHostPort("", config.GetString("http_port"))
	s := http.New(addr, HTTPMux())
	ctx := context.Root
	if err := s.Start(ctx); err != nil {
		zap.L().Fatal("http server start",
			zap.String("addr", addr))
	}
	defer func() {
		if err := s.Stop(ctx); err != nil {
			zap.L().Error("http server stop",
				zap.String("addr", addr))
		}
	}()

}

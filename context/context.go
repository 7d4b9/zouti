package context

import (
	"context"
	"os"
	"os/signal"
)

var Root context.Context

func init() {
	s := make(chan os.Signal)
	signal.Notify(s, os.Interrupt)
	go func() {
		defer close(s)
		var cancel func()
		Root, cancel = context.WithCancel(context.Background())
		<-s
		cancel()
	}()
}

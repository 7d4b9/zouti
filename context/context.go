package context

import (
	"context"
	"os"
	"os/signal"
)

// CancelOnSigInterrupt is a context that will expire after receiving signal SIGINT.
var CancelOnSigInterrupt context.Context

func init() {
	s := make(chan os.Signal)
	signal.Notify(s, os.Interrupt)
	go func() {
		defer close(s)
		var cancel func()
		CancelOnSigInterrupt, cancel = context.WithCancel(context.Background())
		<-s
		cancel()
	}()
}

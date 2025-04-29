package shutdown

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	once              sync.Once
	shutdownBroadcast = make(chan struct{})
)

// Init initializes the shutdown listener.
// It should be called once, typically in main().
func Init(releaseFunc func()) {
	once.Do(func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-c
			releaseFunc()
			close(shutdownBroadcast)
		}()
	})
}

// Done returns a channel that's closed when a shutdown signal is received.
func Done() <-chan struct{} {
	return shutdownBroadcast
}

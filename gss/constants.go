package gss

import (
	"os"
	"syscall"
)

var (
	// DefaultSignalsToWatch The default signals which will trigger the server shutdown.
	DefaultSignalsToWatch = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
)

const (
	// DefaultAddr The default server address.
	DefaultAddr = ":8080"
	// DefaultShutdownTimeoutSeconds the default timeout which the module will wait for
	// connections to drain. -1 means that the module will wait until all connections
	// are drained.
	DefaultShutdownTimeoutSeconds = -1
)

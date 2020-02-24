package gss

import "os"

// Settings the settings which will be used when handling the
// server shutdown.
type Settings struct {
	Addr                   string
	SignalsToWatch         []os.Signal
	ShutdownChannel        chan bool
	ShutdownTimeoutSeconds *int64
}

func (s *Settings) getAddr() string {
	if len(s.Addr) > 0 {
		return s.Addr
	}

	return DefaultAddr
}

func (s *Settings) getSignalsToWatch() []os.Signal {
	if len(s.SignalsToWatch) > 0 {
		return s.SignalsToWatch
	}

	return DefaultSignalsToWatch
}

func (s *Settings) getShutdownTimeoutSeconds() int64 {
	if s.ShutdownTimeoutSeconds != nil {
		return *s.ShutdownTimeoutSeconds
	}

	return DefaultShutdownTimeoutSeconds
}

func (s *Settings) getShutdownChannel() chan bool {
	if s.ShutdownChannel == nil {
		s.ShutdownChannel = make(chan bool)
	}

	return s.ShutdownChannel
}

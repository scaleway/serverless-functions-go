package framework

import (
	"fmt"
	"time"
)

type Option func(*Server)

type Server struct {
	port        string
	coldStart   time.Duration
	baseLatency time.Duration
}

// WithPort can be used to select a port to run the test server, if no port given the server will use an available port given by system and
// print it to stdout.
func WithPort(port int) Option {
	return func(s *Server) {
		s.port = fmt.Sprintf("%d", port)
	}
}

// WithColdStart can be used to simulate cold start. Cold start is the time between the functon call and invocation, this can be meaningful
// for specific clients that have constraints for timeouts.
func WithColdStart(coldStart time.Duration) Option {
	return func(s *Server) {
		s.coldStart = coldStart
	}
}

// With base latency can be used to add time before your function respond to the request. This can be used for example if your request have
// to process some data or to negotiate slow connexions before it responds.
func WithBaseLatency(baseLatency time.Duration) Option {
	return func(s *Server) {
		s.baseLatency = baseLatency
	}
}

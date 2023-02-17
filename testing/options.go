package testing

import (
	"fmt"
)

type Option func(*Server)

type Server struct {
	port string
}

// WithPort can be used to select a port to run the test server, if no port given the server will use an available port given by system and
// print it to stdout.
func WithPort(port int) Option {
	return func(s *Server) {
		s.port = fmt.Sprintf("%d", port)
	}
}

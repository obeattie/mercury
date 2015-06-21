package service

import (
	"github.com/obeattie/mercury/transport"
)

type Config struct {
	Name        string
	Description string
	// Transport specifies a transport to run the server on. If none is specified, a Rabbit transport is used.
	Transport transport.Transport
}

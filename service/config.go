package service

import (
	"github.com/mondough/mercury/transport"
)

type Config struct {
	Name        string
	Description string
	// Transport specifies a transport to run the server on. If none is specified, a mock transport is used.
	Transport transport.Transport
}

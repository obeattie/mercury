package transport

import (
	"sync"

	"github.com/obeattie/typhon/mock"
	ttrans "github.com/obeattie/typhon/transport"
)

type Transport ttrans.Transport

var (
	defaultTransport  Transport
	defaultTransportM sync.RWMutex
)

func init() {
	SetDefaultTransport(mock.NewTransport())
}

// DefaultTransport returns the global default transport, over which servers and clients should run by default
func DefaultTransport() Transport {
	defaultTransportM.RLock()
	defer defaultTransportM.RUnlock()
	return defaultTransport
}

// SetDefaultTransport replaces the global default transport. When replacing, it does not close the prior transport.
func SetDefaultTransport(t Transport) {
	defaultTransportM.Lock()
	defer defaultTransportM.Unlock()
	defaultTransport = t
}

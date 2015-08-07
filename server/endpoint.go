package server

import (
	"runtime"

	log "github.com/cihub/seelog"
	terrors "github.com/mondough/typhon/errors"
	tmsg "github.com/mondough/typhon/message"

	"github.com/mondough/mercury"
	"github.com/mondough/mercury/marshaling"
)

type Handler func(req mercury.Request) (mercury.Response, error)

// An Endpoint represents a handler function bound to a particular endpoint name.
type Endpoint struct {
	// Name is the Endpoint's unique name, and is used to route requests to it.
	Name string
	// Handler is a function to be invoked upon receiving a request, to generate a response.
	Handler Handler
	// Request is a "template" object for the Endpoint's request format.
	Request interface{}
	// Response is a "template" object for the Endpoint's response format.
	Response interface{}
}

func (e Endpoint) unmarshaler(req mercury.Request) tmsg.Unmarshaler {
	return marshaling.Unmarshaler(req.Headers()[marshaling.ContentTypeHeader], e.Request)
}

// Handle takes an inbound Request, unmarshals it, dispatches it to the handler, and serialises the result as a
// Response. Note that the response may be nil.
func (e Endpoint) Handle(req mercury.Request) (rsp mercury.Response, err error) {
	// Unmarshal the request body (unless there already is one)
	if req.Body() == nil && e.Request != nil {
		if um := e.unmarshaler(req); um != nil {
			if err_ := terrors.Wrap(um.UnmarshalPayload(req)); err_ != nil {
				log.Warnf("[Mercury:Server] Cannot unmarshal request payload: %v", err)
				err_.Code = terrors.ErrBadRequest
				rsp, err = nil, err_
				return
			}
		}
	}

	defer func() {
		if v := recover(); v != nil {
			traceVerbose := make([]byte, 1024)
			runtime.Stack(traceVerbose, true)
			log.Criticalf("[Mercury:Server] Recovered from handler panic for request %s:\n%v\n%s", req.Id(), v, traceVerbose)
			rsp, err = nil, terrors.InternalService("Unhandled exception in endpoint handler")
		}
	}()
	rsp, err = e.Handler(req)
	return
}

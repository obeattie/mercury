package server

import (
	log "github.com/cihub/seelog"
	"github.com/golang/protobuf/proto"

	"github.com/obeattie/mercury"
	terrors "github.com/obeattie/typhon/errors"
	tmsg "github.com/obeattie/typhon/message"
)

type Handler func(req mercury.Request) (mercury.Response, error)

// An Endpoint represents a protobuf handler function bound to a particular endpoint name.
type Endpoint struct {
	// Name is the Endpoint's unique name, and is used to route requests to it.
	Name string
	// Handler is a function to be invoked upon receiving a request. Its resulting protobuf Message is serialised as a
	// Response.
	Handler Handler
	// Request is a "template" object for the Endpoint's request format.
	Request proto.Message
	// Response is a "template" object for the Endpoint's response format.
	Response proto.Message
}

// Handle takes an inbound Request, unmarshals it as a protobuf, dispatches it to the handler, and serialises the
// result as a Response. Note that the response may be nil.
func (e Endpoint) Handle(req mercury.Request) mercury.Response {
	// Unmarshal the request body (unless there already is one)
	if req.Body() == nil && e.Request != nil {
		if err := terrors.Wrap(tmsg.ProtoUnmarshaler(e.Request).UnmarshalPayload(req)); err != nil {
			err.Code = terrors.ErrBadRequest
			return e.errorResponse(req, err)
		}
	}

	rsp, rspErr := e.Handler(req)
	if rspErr != nil {
		return e.errorResponse(req, terrors.Wrap(rspErr))
	}
	return rsp
}

// response creates a new correlated Response
func (e Endpoint) response(req mercury.Request) mercury.Response {
	rsp := mercury.NewResponse()
	rsp.SetId(req.Id())
	rsp.SetService(req.Service())
	rsp.SetEndpoint(e.Name)
	return rsp
}

func (e Endpoint) errorResponse(req mercury.Request, err *terrors.Error) mercury.Response {
	log.Debugf("[Mercury:Endpoint] Responding to %s with error: %s", req.Id(), err.Error())
	rsp := e.response(req)
	rsp.SetIsError(true)
	rsp.SetBody(terrors.Marshal(err))
	if err := tmsg.ProtoMarshaler().MarshalBody(rsp); err != nil {
		log.Errorf("[Mercury:Endpoint] Failed to marshal error response for %s: %v", req.Id(), err)
		return nil // Not much we can do here :(
	}
	return rsp
}

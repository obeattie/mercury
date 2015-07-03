package server

import (
	"github.com/golang/protobuf/proto"

	terrors "github.com/mondough/typhon/errors"
	tmsg "github.com/mondough/typhon/message"
	"github.com/obeattie/mercury"
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
func (e Endpoint) Handle(req mercury.Request) (mercury.Response, error) {
	var err *terrors.Error

	// Unmarshal the request body (unless there already is one)
	if req.Body() == nil && e.Request != nil {
		if err = terrors.Wrap(tmsg.ProtoUnmarshaler(e.Request).UnmarshalPayload(req)); err != nil {
			err.Code = terrors.ErrBadRequest
			return nil, err
		}
	}

	return e.Handler(req)
}

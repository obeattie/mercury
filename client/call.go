package client

import (
	"golang.org/x/net/context"

	terrors "github.com/mondough/typhon/errors"
	tmsg "github.com/mondough/typhon/message"
	"github.com/obeattie/mercury"
	"github.com/obeattie/mercury/marshaling"
)

// A Call is a convenient way to form a Request for an RPC call.
type Call struct {
	// Uid represents a unique identifier for this call; it is used.
	Uid string
	// Service to receive the call.
	Service string
	// Endpoint of the receiving service.
	Endpoint string
	// Body will be serialised to form the Payload of the request.
	Body interface{}
	// Headers to send on the request (these may be augmented by the client).
	Headers map[string]string
	// Response is a protocol into which the response's Payload should be unmarshaled.
	Response interface{}
	// Context is a context for the request. This should nearly always be the parent request (if any).
	Context context.Context
}

func (c Call) marshaler() tmsg.Marshaler {
	result := tmsg.Marshaler(nil)
	if c.Headers != nil && c.Headers[marshaling.ContentTypeHeader] != "" {
		result = marshaling.Marshaler(c.Headers[marshaling.ContentTypeHeader])
	}
	if result == nil {
		result = tmsg.ProtoMarshaler()
	}
	return result
}

// Request yields a Request formed from this Call
func (c Call) Request() (mercury.Request, error) {
	req := mercury.NewRequest()
	req.SetService(c.Service)
	req.SetEndpoint(c.Endpoint)
	req.SetHeaders(c.Headers)
	if c.Context != nil {
		req.SetContext(c.Context)
	}
	if c.Body != nil {
		req.SetBody(c.Body)
		if err := c.marshaler().MarshalBody(req); err != nil {
			terr := terrors.Wrap(err)
			terr.Code = terrors.ErrBadRequest
			return nil, terr
		}
	}
	return req, nil
}

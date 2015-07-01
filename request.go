package mercury

import (
	"sync"
	"time"

	log "github.com/cihub/seelog"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	tmsg "github.com/mondough/typhon/message"
)

const (
	errHeader         = "Content-Error"
	ContentTypeHeader = "Content-Type"
	acceptHeader      = "Accept"
	ProtoContentType  = tmsg.ProtoContentType
	JSONContentType   = tmsg.JSONContentType
)

// A Request is a representation of an RPC call (inbound or outbound). It extends Typhon's Request to provide a
// Context, and also helpers for constructing a response.
type Request interface {
	tmsg.Request
	context.Context

	// Response constructs a response to this request, with the (optional) given body proto. The response will share
	// the request's ID, and be destined for the originator.
	Response(body proto.Message) Response
	// A Context for the Request.
	Context() context.Context
	// SetContext replaces the Request's Context.
	SetContext(ctx context.Context)
}

func responseFromRequest(req Request, body proto.Message) Response {
	rsp := NewResponse()
	rsp.SetId(req.Id())
	rsp.SetService(req.Service())
	rsp.SetEndpoint(req.Endpoint())
	if body != nil {
		rsp.SetBody(body)
		var err error
		switch req.Headers()[acceptHeader] {
		case JSONContentType:
			err = tmsg.JSONMarshaler().MarshalBody(rsp)
		default:
			err = tmsg.ProtoMarshaler().MarshalBody(rsp)
		}
		if err != nil {
			log.Errorf("[Mercury] Failed to marshal response %s: %v", rsp.Id(), err)
			return nil
		}
	}
	return rsp
}

type request struct {
	sync.RWMutex
	tmsg.Request
	ctx context.Context
}

func (r *request) Response(body proto.Message) Response {
	return responseFromRequest(r, body)
}

func (r *request) Context() context.Context {
	r.RLock()
	defer r.RUnlock()
	return r.ctx
}

func (r *request) SetContext(ctx context.Context) {
	r.Lock()
	defer r.Unlock()
	r.ctx = ctx
}

func (r *request) Copy() tmsg.Request {
	r.RLock()
	defer r.RUnlock()
	return &request{
		Request: r.Request.Copy(),
		ctx:     r.ctx,
	}
}

// Context implementation

func (r *request) Deadline() (time.Time, bool) {
	return r.Context().Deadline()
}

func (r *request) Done() <-chan struct{} {
	return r.Context().Done()
}

func (r *request) Err() error {
	return r.Context().Err()
}

func (r *request) Value(key interface{}) interface{} {
	return r.Context().Value(key)
}

func NewRequest() Request {
	return FromTyphonRequest(tmsg.NewRequest())
}

func FromTyphonRequest(req tmsg.Request) Request {
	return &request{
		Request: req,
		ctx:     context.TODO(),
	}
}

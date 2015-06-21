package requesttree

import (
	"golang.org/x/net/context"

	"github.com/obeattie/mercury"
	"github.com/obeattie/mercury/client"
	"github.com/obeattie/mercury/server"
)

const (
	parentIdHeader = "Parent-Request-ID"
	reqIdCtxKey    = "Request-ID"
)

var sharedMiddleware = &rtm{}

type ParentRequestIdMiddleware interface {
	server.ServerMiddleware
	client.ClientMiddleware
}

type rtm struct{}

func (m *rtm) ProcessClientRequest(req mercury.Request) mercury.Request {
	if req.Headers()[parentIdHeader] == "" { // Don't overwrite an exiting header
		if parentId, ok := req.Context().Value(reqIdCtxKey).(string); ok && parentId != "" {
			req.SetHeader(parentIdHeader, parentId)
		}
	}
	return req
}

func (m *rtm) ProcessClientResponse(rsp mercury.Response, ctx context.Context) mercury.Response {
	return rsp
}

func (m *rtm) ProcessServerRequest(req mercury.Request) mercury.Request {
	req.SetContext(context.WithValue(req.Context(), reqIdCtxKey, req.Id()))
	if v := req.Headers()[parentIdHeader]; v != "" {
		req.SetContext(context.WithValue(req.Context(), parentIdCtxKey, v))
	}
	return req
}

func (m *rtm) ProcessServerResponse(rsp mercury.Response, ctx context.Context) mercury.Response {
	if v, ok := ctx.Value(parentIdCtxKey).(string); ok && v != "" {
		rsp.SetHeader(parentIdHeader, v)
	}
	return rsp
}

func Middleware() ParentRequestIdMiddleware {
	return sharedMiddleware
}

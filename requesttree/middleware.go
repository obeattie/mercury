package requesttree

import (
	"golang.org/x/net/context"

	"github.com/mondough/mercury"
	terrors "github.com/mondough/typhon/errors"
)

const (
	parentIdHeader = "Parent-Request-ID"
	reqIdCtxKey    = "Request-ID"
)

type requestTreeMiddleware struct{}

func (m requestTreeMiddleware) ProcessClientRequest(req mercury.Request) mercury.Request {
	if req.Headers()[parentIdHeader] == "" { // Don't overwrite an exiting header
		if parentId, ok := req.Context().Value(reqIdCtxKey).(string); ok && parentId != "" {
			req.SetHeader(parentIdHeader, parentId)
		}
	}

	// Pass through the current service and endpoint as the origin of this request
	if svc, ok := req.Value("Service").(string); ok {
		req.SetHeader("Origin-Service", svc)
	}
	if ept, ok := req.Value("Endpoint").(string); ok {
		req.SetHeader("Origin-Endpoint", ept)
	}

	return req
}

func (m requestTreeMiddleware) ProcessClientResponse(rsp mercury.Response, ctx context.Context) mercury.Response {
	return rsp
}

func (m requestTreeMiddleware) ProcessClientError(err *terrors.Error, ctx context.Context) {}

func (m requestTreeMiddleware) ProcessServerRequest(req mercury.Request) (mercury.Request, mercury.Response) {
	req.SetContext(context.WithValue(req.Context(), reqIdCtxKey, req.Id()))
	if v := req.Headers()[parentIdHeader]; v != "" {
		req.SetContext(context.WithValue(req.Context(), parentIdCtxKey, v))
	}

	// Set the current service and endpoint into the context
	req.SetContext(context.WithValue(req.Context(), "Service", req.Service()))
	req.SetContext(context.WithValue(req.Context(), "Endpoint", req.Endpoint()))

	// Set the originator into the context
	req.SetContext(context.WithValue(req.Context(), "Origin-Service", req.Headers()["Origin-Service"]))
	req.SetContext(context.WithValue(req.Context(), "Origin-Endpoint", req.Headers()["Origin-Endpoint"]))

	return req, nil
}

func (m requestTreeMiddleware) ProcessServerResponse(rsp mercury.Response, ctx context.Context) mercury.Response {
	if v, ok := ctx.Value(parentIdCtxKey).(string); ok && v != "" && rsp != nil {
		rsp.SetHeader(parentIdHeader, v)
	}
	return rsp
}

func Middleware() requestTreeMiddleware {
	return requestTreeMiddleware{}
}

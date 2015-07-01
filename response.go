package mercury

import (
	tmsg "github.com/mondough/typhon/message"
)

type Response interface {
	tmsg.Response
	IsError() bool
	SetIsError(v bool)
}

type response struct {
	tmsg.Response
}

func (r *response) IsError() bool {
	return r.Headers()[errHeader] == "1"
}

func (r *response) SetIsError(v bool) {
	r.SetHeader(errHeader, "1")
}

func NewResponse() Response {
	return FromTyphonResponse(tmsg.NewResponse())
}

func FromTyphonResponse(rsp tmsg.Response) Response {
	return &response{
		Response: rsp,
	}
}

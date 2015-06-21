package client

import (
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

// Req sends a synchronous request to a service using a new client, and unmarshals the response into the supplied
// protobuf
func Req(ctx context.Context, service, endpoint string, req, res proto.Message) error {
	if err := NewClient().
		Add(Call{
		Uid:      "1",
		Service:  service,
		Endpoint: endpoint,
		Body:     req,
		Response: res,
		Context:  ctx,
	}).Execute().Errors().Combined(); err != nil {
		return err
	} else {
		// Relevant: http://golang.org/doc/faq#nil_error
		return nil
	}
}

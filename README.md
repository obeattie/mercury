# Mercury

An RPC client/server implementation using [Typhon](https://github.com/obeattie/typhon), intended for building microservices.

## Server

A [`Server`](http://godoc.org/github.com/obeattie/mercury/server) receives RPC requests, routes them to an [`Endpoint`](http://godoc.org/github.com/obeattie/mercury/server#Endpoint), calls a handler function to "do work," and returns a response back to a caller.

### Server middleware

Server middleware offers hooks into request processing for globally altering a server's input or output. They could be used to provide authentication or distributed tracing functionality, for example.

## Client

A [`Client`](http://godoc.org/github.com/obeattie/mercury/client#Client) offers a convenient way, atop a Typhon transport, to make requests to other servers. They co-ordinate the execution of many parallel requests, deal with response and error unmarshaling, and provide convenient ways of dealing with response errors.

### Client middleware

Like server middleware, clients too have hooks for altering outbound requests or inbound responses.

## Service

A [`Service`](http://godoc.org/github.com/obeattie/mercury/service#Service) is a lightweight wrapper around a server, which also sets up some global defaults (for instance, to use the same default client transport as the server).

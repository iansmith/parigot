package main

import (
	greetv1connect "example/greet/proto/gen/greet/greetconnect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
)

// where should this go
type RequestIdT int64
type ModeT int64

// where should this go?
type ParigotCtx struct {
	Mode      ModeT
	RequestId RequestIdT
	Header    string
}

type ParigotError struct {
	Message string
}

func (p *ParigotError) Error() string {
	return p.Message
}

type NotImplementedYet struct {
	*ParigotError
}

func NewNotImpletedYet() error {
	return &NotImplementedYet{
		&ParigotError{Message: "not implemented yet"},
	}
}

type BaseService interface{}

type BaseServiceImpl struct {
}

func (b *BaseServiceImpl) Connect(s string) *BaseServiceImpl {
	//return things from the api of parigot
}

func lib_mainloop() {
	greeter := &GreetServer{}
	mux := http.NewServeMux()
	path, handler := greetv1connect.NewGreetHandler(greeter)
	mux.Handle(path, handler)
	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}

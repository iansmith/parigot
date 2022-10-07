package main

import "example/greet/proto/gen/greet"

type GreetImpl struct {
	*BaseServiceImpl
}

func NewGreet() *GreetImpl {
	return &GreetImpl{
		BaseServiceImpl: &BaseServiceImpl{},
	}
}

func (g *GreetImpl) Start() {
}

func (g *GreetImpl) Greet(p *ParigotCtx, req greet.GreetRequest) (string, error) {
	return "hi there, " + req.Name, nil
}

// creator
func NewGreetImpl() Greet {
	return &GreetImpl{}
}

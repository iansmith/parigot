package main

type NotClearYet struct {
}

func (g *NotClearYet) Greet(_ ParigotCtx, _ string) (string, error) {
	return "", NewNotImpletedYet()
}

// generated code
type Greet interface {
	BaseService
	Greet()
}

func lib_Register(org string, name string, fn func() Greet) {
	// how to implement this?
}

func lib_Connect[T](string, string) Greet {

}

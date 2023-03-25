package main

import (
	"github.com/iansmith/parigot/apiimpl/dom"
	dommsg "github.com/iansmith/parigot/g/msg/dom/v1"
)

const parentId = "paraLoc"

var exitChan = make(chan bool)

func main() {

	svc, err := dom.LocateDOMServer()
	if err != nil {
		panic("unable to get DOMService: " + err.Error())
	}
	resp, err := svc.ElementById(&dommsg.ElementByIdRequest{Id: parentId})
	if err != nil {
		panic("unable to get element by id:" + err.Error())
	}
	req := &dommsg.CreateElementRequest{
		Root:   example(),
		Parent: resp.Elem,
	}
	_, err = svc.CreateElement(req)
	if err != nil {
		panic("failed to create element: " + err.Error())
	}

	AddGlobalEvent(svc)
	_ = <-exitChan
}

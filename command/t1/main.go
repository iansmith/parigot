package main

import (
	"log"

	"github.com/iansmith/parigot/apiimpl/dom"
	dommsg "github.com/iansmith/parigot/g/msg/dom/v1"
	filemsg "github.com/iansmith/parigot/g/msg/file/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"
	queuemsg "github.com/iansmith/parigot/g/msg/queue/v1"
	lib "github.com/iansmith/parigot/lib/go"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const parentId = "paraLoc"

var exitChan = make(chan bool)

func Main() {

	svc, err := dom.LocateDOMServer()
	if err != nil {
		panic("unable to get DOMService: " + err.Error())
	}
	resp, err := svc.ElementById(&dommsg.ElementByIdRequest{Id: parentId})
	if err != nil {
		panic("unable to get element by id:" + err.Error())
	}
	// this value is the PAYLOAD being sent in the QueueMsg... this is just
	// to test that the values get passed through correctly
	payload := &filemsg.CreateRequest{
		Path: "/foo/bar/baz",
	}
	var a anypb.Any
	err = (&a).MarshalFrom(payload)
	if err != nil {
		log.Fatalf("unable to marshal payload: %v", err)
	}
	mod := &queuemsg.QueueMsg{
		Id:           lib.Marshal[protosupportmsg.QueueId](lib.NewQueueId()),
		MsgId:        lib.Marshal[protosupportmsg.QueueMsgId](lib.NewQueueMsgId()),
		ReceiveCount: 0,
		Received:     timestamppb.Now(),
		Sender:       nil,
		Sent:         timestamppb.Now(),
		Payload:      &a,
	}

	req := &dommsg.CreateElementRequest{
		Root:   QueueMsgParent(mod),
		Parent: resp.Elem,
	}
	_, err = svc.CreateElement(req)
	if err != nil {
		panic("failed to create element: " + err.Error())
	}
	elem := QueueMsgViewShort(mod, 0)
	if elem == nil {
		log.Fatalf("unable to create the initial short display")
	}
	setChildReq := &dommsg.SetChildRequest{
		Id:        lib.Unmarshal(mod.Id).String(), // we just grab it from the already built msg object
		ParigotId: elem.ParigotId,
		Child:     []*dommsg.Element{elem},
	}
	_, err = svc.SetChild(setChildReq)
	if err != nil {
		log.Fatalf("unable to set child: %+v=> %v", setChildReq, err)
	}
	AddGlobalEvent(svc)
	_ = <-exitChan
}

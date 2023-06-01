//
// DO NOT EDIT.  This file was machine generated by WCL from ui/testdata/model_test.wcl.
//

package main

import (
	"bytes"
	"fmt"
	"syscall/js"

	apidom "github.com/iansmith/parigot/apiimpl/dom"
	dom "github.com/iansmith/parigot/g/dom/v1"
	dommsg "github.com/iansmith/parigot/g/msg/dom/v1"
	protosupportmsg "github.com/iansmith/parigot/g/msg/protosupport/v1"

	"github.com/iansmith/parigot/g/msg/queue/v1"
	"github.com/iansmith/parigot/lib/go"
	"github.com/iansmith/parigot/ui/parser/builtin"

	"google.golang.org/protobuf/types/known/anypb"
)

// This is necessary so we don't have to refer to builtin.ParigotId which causes
// parsing problems.
func ParigotId[T lib.AllIdPtr](id T) string {
	return builtin.ParigotId(id)
}
func switchShortAndLong(this js.Value, arg []js.Value) any {
	return nil
}

//
// Text Section
//

func sizeInBytes(p *anypb.Any) string {
	var size int
	var result bytes.Buffer
	size = len(p.Value)

	result.WriteString(fmt.Sprint(size))
	result.WriteString(` bytes`)

	return result.String()
}

func shortModel(model *queuemsg.QueueMsg) string {
	var qid string
	var mid string
	var size string
	var result bytes.Buffer
	qid = lib.Unmarshal(model.Id).Short()
	mid = lib.Unmarshal(model.MsgId).Short()
	size =
		fmt.Sprint(sizeInBytes(model.Payload))

	result.WriteString(`QueueId: `)
	result.WriteString(fmt.Sprint(qid))
	result.WriteString(`:Message Id: `)
	result.WriteString(fmt.Sprint(mid))
	result.WriteString(`:`)
	result.WriteString(fmt.Sprint(size))

	return result.String()
}

//
// Doc Section
//

func QueueMsgParent(model *queuemsg.QueueMsg) *dommsg.Element {

	// number is 0
	result :=
		&dommsg.Element{Tag: &dommsg.Tag{
			Name: "div", Id: ParigotId(model.Id),
		}, // end of tag with name,id,class

			// tag part ended
			ParigotId: lib.Marshal[protosupportmsg.ElementId](lib.NewElementId()),
		} // end of doc element (with no Children)

	return result
}

func QueueMsgView(model *queuemsg.QueueMsg, a int64) *dommsg.Element {

	// number is 1
	n1 :=
		&dommsg.Element{Tag: &dommsg.Tag{
			Name: "h5",
		}, // end of tag with name,id,class

			// tag part ended
			ParigotId: lib.Marshal[protosupportmsg.ElementId](lib.NewElementId()),

			Text: "" + `some text before ` + fmt.Sprint(model.MsgId.Id) + ` and some text after`,
		} // end of doc element (with no Children)

	// number is 2
	n2 :=
		&dommsg.Element{Tag: &dommsg.Tag{
			Name: "h6",
		}, // end of tag with name,id,class

			// tag part ended
			ParigotId: lib.Marshal[protosupportmsg.ElementId](lib.NewElementId()),

			Text: "" + `Sent: ` + fmt.Sprint(model.Sent),
		} // end of doc element (with no Children)

	// number is 3
	n3 :=
		&dommsg.Element{Tag: &dommsg.Tag{
			Name: "h6",
		}, // end of tag with name,id,class

			// tag part ended
			ParigotId: lib.Marshal[protosupportmsg.ElementId](lib.NewElementId()),

			Text: "" + `Size:` + fmt.Sprint(sizeInBytes(model.Payload)) + ` bytes`,
		} // end of doc element (with no Children)

	// number is 0
	result :=
		&dommsg.Element{Tag: &dommsg.Tag{
			Name: "h4",
		}, // end of tag with name,id,class

			// tag part ended
			ParigotId: lib.Marshal[protosupportmsg.ElementId](lib.NewElementId()),

			Text: "" + `Message ` + fmt.Sprint(model.MsgId.Id) + ` from Queue ` + fmt.Sprint(model.Id),
			Child: []*dommsg.Element{
				n1,
				n2,
				n3,
			}, // end of children
		}

	return result
}

func QueueMsgViewShort(model *queuemsg.QueueMsg, a int64) *dommsg.Element {

	// number is 0
	result :=
		&dommsg.Element{Tag: &dommsg.Tag{
			Name: "h4",
		}, // end of tag with name,id,class

			// tag part ended
			ParigotId: lib.Marshal[protosupportmsg.ElementId](lib.NewElementId()),

			Text: "" + fmt.Sprint(shortModel(model)),
		} // end of doc element (with no Children)

	return result
}

//
// Event Section
//

func AddEventQueueMsg(model *queuemsg.QueueMsg, svc dom.DOMService) {
	selector := "#" + ParigotId(model.MsgId)
	fn := func(this js.Value, arg []js.Value) any {
		return switchShortAndLong(this, append(arg, js.ValueOf(model), js.ValueOf(selector), js.ValueOf(svc.(*apidom.DOMServer).ServerId())))
	}
	svc.(*apidom.DOMServer).AddEvent("#"+ParigotId(model.MsgId), "click", fn)
}
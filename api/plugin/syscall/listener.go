package syscall

import (
	"log"
	"reflect"
	"time"

	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/g/syscall/v1"
)

type ListenerKind int

const (
	UnknownKind ListenerKind = iota
	Timeout
	MethodCall
	Exit
)

// Listener is a type that is used when we are running
type Listener interface {
	Case() []reflect.SelectCase
	// Handle sets up the return value whose pointer is
	// given in the params.  The first parameter is the value
	// sent through the channel.  The second parameter is the choice
	// number of the case selected--this probably only matters
	// a methodRequestListener.
	Handle(reflect.Value, int, *syscall.ReadOneResponse)
}

type methodCallListener struct {
	req *syscall.ReadOneRequest
}

func newMethodCallListener(req *syscall.ReadOneRequest) *methodCallListener {
	mcl := &methodCallListener{
		req: req,
	}
	return mcl
}

func (m *methodCallListener) Case() []reflect.SelectCase {
	cases := make([]reflect.SelectCase, len(m.req.GetCall()))
	for i, pair := range m.req.Call {
		svc := id.UnmarshalServiceId(pair.ServiceId)
		meth := id.UnmarshalMethodId(pair.MethodId)
		combo := makeSidMidCombo(svc, meth)
		ch := pairIdToChannel[combo]
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	return cases
}

func (m *methodCallListener) Handle(value reflect.Value, chosen int, resp *syscall.ReadOneResponse) {
	pair := m.req.Call[chosen]

	resp.Call = &syscall.ServiceMethodCall{}
	resp.Call.ServiceId = pair.ServiceId
	resp.Call.MethodId = pair.MethodId
	resp.CallId = value.Interface().(CallInfo).cid.Marshal()
	resp.Param = value.Interface().(CallInfo).param

	resp.Timeout = false
	resp.Exit = false
}

//
// ExitListener
//

type ExitListener struct {
	ch chan int32
}

func NewExitListener(ch chan int32) *ExitListener {
	return &ExitListener{ch}
}

func (e *ExitListener) Case() []reflect.SelectCase {
	if e.ch == nil {
		return nil
	}
	return []reflect.SelectCase{
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(e.ch),
		},
	}
}

// Handle only needs to look at the value received, since it has
// an exit code in it.
func (e *ExitListener) Handle(value reflect.Value, _ int, resp *syscall.ReadOneResponse) {
	if !value.CanInt() {
		log.Printf("unable to understand value provided to ExitListener Handle (%d)", value.Kind())
		return
	}
	resp.Call = nil
	resp.CallId = nil
	resp.Param = nil

	resp.Timeout = false
	resp.Exit = true
}

//
// Timeout
//

type timeoutListener struct {
	timeout int32
}

func newTimeoutListener(timeoutInMillis int32) *timeoutListener {
	return &timeoutListener{
		timeout: timeoutInMillis,
	}
}

func (t *timeoutListener) Case() []reflect.SelectCase {
	if t.timeout >= 0 {
		ch := time.After(time.Duration(t.timeout) * time.Millisecond)

		return []reflect.SelectCase{
			{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ch),
			},
		}
	}
	return nil
}

// Handle knows it's a timeout, that's all that matters here.
// We return the response with the boolean for timeout set.
func (t *timeoutListener) Handle(_ any, _ int, resp *syscall.ReadOneResponse) {
	resp.Call = nil
	resp.CallId = nil
	resp.Param = nil

	resp.Timeout = true
	resp.Exit = false
}

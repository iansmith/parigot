package kernel

import (
	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/proto"
)

// meta-note: All the system calls are called on the kernel.
// meta-note: In some cases the kernel simply forwards the request and response
// meta-note: to another system, i.e. Export() forwarded to starter.
// meta-note: When the kernel just forwards something, the request and response
// meta-note: pointers are passed and when a method on a subsystem has
// meta-note: the same name as a system call, but not the req/resp pair, then the
// meta-note: kernel processes part of the request and the subsystem a different part.

// K is the kernel.  It is initialized by kernel.InitSingle or other
// kernel.InitXXX().  Each of these configure the kernel in different ways.
// This value should not set more than once.
var K Kernel

// Kernel is... well... the kernel.
type Kernel interface {
	// CancelRead does nothing if the kernel is not waiting. If the kernel is
	// waiting, this method causes that to stop and the kernel returnns a timeout
	// from it's ReadOne() method.
	CancelRead()

	// SetApproach should be called once, at startup, to indicate what type of
	// deployment you are using.  Each type of deployment has exactly one approach
	// that does message send, receive, and finish.  Usually the Nameserver and
	// Starter also need to be coordinated to make the approach work.
	SetApproach(GeneralReceiver, GeneralReceiver, Nameserver, Starter) syscall.KernelErr

	// AddRecevier is the generic version of a receiver.  This is usually the most
	// useful if you want listen for an external network protocol.
	AddReceiver(GeneralReceiver)

	// Register is used to notify the kernel that a given service
	// should be assigned a service id.  Note that this may reach multiple
	// parts of the kernel based on the Registrar interface.
	Register(req *syscall.RegisterRequest, resp *syscall.RegisterResponse) syscall.KernelErr

	// Dispatch is used to send a call to a remote machine.  If this
	// returns a kernel error it is because the dispatch call itself could
	// not be made, not that the dispatch worked ok and an error was returned
	// by the remote code.
	Dispatch(req *syscall.DispatchRequest, resp *syscall.DispatchResponse) syscall.KernelErr

	// Launch logically causes a process to wait for all its dependencies to
	// be ready.  In practice, it returns immediately and then finishes the
	// process later.
	Launch(req *syscall.LaunchRequest, resp *syscall.LaunchResponse) syscall.KernelErr

	// BindMethod creates the mapping from the name of a method on a service to the
	// method id (and corresponding service id) that one can use to call
	// the method.
	BindMethod(req *syscall.BindMethodRequest, resp *syscall.BindMethodResponse) syscall.KernelErr

	// Export connects a particular service id to a named
	Export(req *syscall.ExportRequest, resp *syscall.ExportResponse) syscall.KernelErr

	// ReadOne returns the response to a (blocking) ReadOne request.
	// This queries various parts of the system to gather information
	// that may be needed for the response.
	ReadOne(req *syscall.ReadOneRequest, resp *syscall.ReadOneResponse) syscall.KernelErr

	// Require declares a dependency(ies) between the source and destination.
	// Destination must be started before source.
	Require(req *syscall.RequireRequest, resp *syscall.RequireResponse) syscall.KernelErr

	// Locate is the constructor for the types in parigot. It takes the name of
	// an interface (like "foo.v1.Bar") a returns a service id that implements
	// that service.
	Locate(req *syscall.LocateRequest, resp *syscall.LocateResponse) syscall.KernelErr

	// ReturnValue is used to finish a previous Dispatch call.  This is where the
	// original caller will get his call completed.
	ReturnValue(req *syscall.ReturnValueRequest, resp *syscall.ReturnValueResponse) syscall.KernelErr

	// Exit is a call that can be called to exit a single program
	// or the whole network of services.
	Exit(req *syscall.ExitRequest, resp *syscall.ExitResponse) syscall.KernelErr

	// Nameserver gets the nameserver for the kernel.  This
	// does not lock.
	Nameserver() Nameserver
}

// Maximum number of different channels to select on
const maxCases = 32

// what is the maximum time we will wait inside a read one
const maxWait = 10000 //millis

// Registrar is an interface that gets notified when any
// services registers. This effectively allows the implementor
// to become of aware of all services that are created and their
// corresponding host.
type Registrar interface {
	Register(hid id.HostId, sid id.ServiceId, debugName string) syscall.KernelErr
}

// Binder is an interface that gets notified when any
// method is bound.  It connects the name of the method with the service
// and method id.
type Binder interface {
	Bind(hid id.HostId, sid id.ServiceId, mid id.MethodId, methodName string) syscall.KernelErr
}

// Starter is the type that handles the machinery of getting the
// the services started up in the right order. "right" here means that
// a service doesn't start before all of it's dependencies
// are ready.  Further, it will generally refuse attempts to locate, or "find"
// dependencies that were not previously declared.
type Starter interface {
	Registrar
	Binder
	// Require is used to indicate that a given service cannot run until all the given
	// services are already running.
	Require(*syscall.RequireRequest, *syscall.RequireResponse) syscall.KernelErr
	// Ready returns the next service id that can run and
	// the number of services left to consider.  If the returned
	// service id is the zero value, then no service is ready to run.
	Ready() (launchCompleteBundle, int)
	// Export declares that the given service exports one or more interfaces,
	// as strings like "foo.v1.Bar".
	Export(*syscall.ExportRequest, *syscall.ExportResponse) syscall.KernelErr
	// Launch is called by a service indicating that all its preliminaries are
	// complete and it is waiting for the starter the tell
	// it to run.
	Launch(sid id.ServiceId, cid id.CallId, hid id.HostId, mid id.MethodId) syscall.KernelErr
	// This is the way that a named interface gets turned into a service
	// that implements that
	Locate(req *syscall.LocateRequest, resp *syscall.LocateResponse) syscall.KernelErr
}

type GeneralReceiver interface {
	Ch() chan *syscall.ReadOneResponse
	TimeoutInMillis() int
}

// type GeneralSender interface {
// 	Send(any) syscall.KernelErr
// 	TimeoutInMillis() int
// }

// Receiver only
type HttpConnector struct {
}

// Sender and Finisher but not receiver
type HttpProxy struct {
}

// Sender only
type KernelLogger struct {
}

type KLog interface {
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
}

type Nameserver interface {
	Registrar
	Binder

	// AllHosts enumerates all the hosts that are known.
	AllHosts() []id.HostId

	// AddHosts adds a new host to the set of known hosts. Adding
	// a host multiple times is allowed.
	AddHost(id.HostId)

	// FindHost returns the host for a service id. It returns the
	// HostId zero value if this fails.
	FindHost(id.ServiceId) id.HostId

	// FindHosChan returns the channel one can write on to send a message
	// to that host. It return nil if this fails.
	FindHostChan(id.HostId) chan<- proto.Message

	// FindMethod returns the name of a method, given host, service, and method id.
	FindMethod(id.HostId, id.ServiceId, id.MethodId) string

	// In() requests a chan that can be used to read requests
	// from the network.
	In() chan proto.Message
}

// fqName is a fully qualified name of a service, analagous to
// syscall.FullyQualifiedName.
type fqName struct {
	pkg, name string
}

// MakeSidMidCombo is a utility for construction of a key (string) that is
// derived from the sid and mid given.
func MakeSidMidCombo(sid id.ServiceId, mid id.MethodId) string {
	return sid.String() + "," + mid.String()
}

// launchCompleteBundle is returned by the starter when a particular
// service is ready to run.  The service is running even if the
// hasCycle value is true. The caller needs to take action if it
// wishes to have cyclic dependencies stop the program.
// In the the case where there is nothing launched, sid will be
// the zero value.
type launchCompleteBundle struct {
	hid      id.HostId
	cid      id.CallId
	sid      id.ServiceId
	hasCycle bool
}

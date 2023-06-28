package syscall

import (
	"context"

	"github.com/iansmith/parigot/apishared/id"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/types/known/anypb"
)

// Service is the logical representation of a service. This is
// used internally and is not intended for external use.
type Service interface {
	Id() id.ServiceId
	// Name returns a human readable name of the service.
	Name() string
	// Package returns the package name (not the proto package name)
	// of the service.
	Package() string
	// Short returns a nice-to-read version of the service's id.
	Short() string
	// String returns the full id of the service.
	String() string
	// RunRequested returns true if the service has requested
	// to run, but its dependencies are not yet satisfied. Once
	// they are met, the service can start and the Started()
	// method will return true.
	RunRequested() bool
	// Started returns true if the service has started.
	Started() bool
	// Exported returns true if some service provider has said that
	// they implement this service.
	Exported() bool
	// Method returns all the pairs of MethodName and MethodId
	// for a service known to the SyscallData.  You provide the
	// service to this method to know which set of pairs you want.
	Method() []*syscall.MethodBinding
	// AddMethod is called by the syscall bind method to add a
	// given name/id pair to this service.
	AddMethod(string, id.MethodId)
	//Run is badly named. This really means "block until everything
	//I need is ready."
	Run(context.Context) syscall.KernelErr
}

// SyscallData is the interface used by the kernel methods
// (syscallhost.go) to get information about the status of
// a startup sequence.

type SyscallData interface {
	//ServiceByName looks up a service and returns it based on the
	//values package_ and name.  If this returns nil, the service could
	//not be found.
	ServiceByName(ctx context.Context, package_, name string) Service
	//ServiceById looks up a service and returns it based on the
	//value sid.  If this returns nil, the service could
	//not be found.
	ServiceById(ctx context.Context, sid id.ServiceId) Service
	//ServiceByIdString looks up a service based on the printed representation
	//of the service id.  If the service cannot be found ServiceByIdString
	//returns nil.
	ServiceByIdString(ctx context.Context, str string) Service
	// SetService puts a service into SyscallData.  This should only be
	// called once for each package_ and name pair. It returns the
	// ServiceId for the service named, creating a new one if necessary.
	// The client flag should be set to true only when the requesting
	// party is a client.  All services should pass false here.  This
	// flag effectively means that the requester (package_,name) does not
	// to export their service to be ready to run.
	// If the bool result is false, then the pair already existed and
	// we made no changes to it.
	SetService(ctx context.Context, package_, name string, client bool) (Service, bool)
	// Export finds a service by the given sid and then marks that
	// service as being exported. This function returns nil if
	// there is no such service.
	Export(ctx context.Context, svc id.ServiceId) Service
	// Import introduces a dendency between the sourge and dest
	// services. Thus,  dest must be running before source can run.
	// This function returns a kernel error in two primary cases.
	// 1. one of the src or destination could not be found.  2. The
	// newly introduced edge would create a cycle.
	Import(ctx context.Context, src, dest id.ServiceId) syscall.KernelErr
	// Launch blocks the caller until all the prerequistes have been
	// launched.  It returns false if it returned because of
	// a timeout or the service id cannot be found, otherwise true.
	Launch(context.Context, id.ServiceId) syscall.KernelErr
	// PathExists returns true if there is a sequence of dependency
	// graph vertices that eventually leads from source to target.
	PathExists(ctx context.Context, source, target string) bool
}

// CallMatcher is an internal data structure object that
// connects calls to Dispatch (the call) with the response
// which are created by ReturnValue requests.
type CallMatcher interface {
	// Response is called when a return value is
	// being processed. Any value that
	// is returned is NOT from the execution but from
	// the Response call itself.  Be aware that the
	// Response call is likely to be from a different
	// host than the original Dispatch call.
	Response(cid id.CallId, a *anypb.Any, err int32) syscall.KernelErr
	// Dispatch creates the necessary entries to handle
	// a future call to Response.  The value returned is
	// related to the Dispatch itself, it is not related
	// to the execution of the call being registered.
	Dispatch(hid id.HostId, cid id.CallId) syscall.KernelErr
	// Ready returns a resolved call or nil if no Promises are
	// resolved for the given host.
	Ready(hid id.HostId) (*syscall.ResolvedCall, syscall.KernelErr)
}

// HostFinder returns information about a host in the format used
// by the syscall struct. It is convention to use the fully qualified
// name of the service for the name.
type HostFinder interface {
	// FindByName finds the correct host by the name field.
	// If the name cannot be found, it returns nil.
	FindByName(name string) *hostInfo
	// FindById finds the correct host by the id field.
	// If the id cannot be found it returns nil.
	FindById(id id.HostId) *hostInfo
	// AddHost is used to add a record to the set of hosts
	// that are know. This call will panic if either the
	// name or id is not set.
	AddHost(name string, hid id.HostId) syscall.KernelErr
}

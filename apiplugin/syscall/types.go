package main

import (
	"context"

	"github.com/iansmith/parigot/apishared/id"
)

type Service interface {
	Id() id.ServiceId
	Name() string
	Package() string
	Short() string
	String() string
	RunRequested() bool
	Started() bool
	Exported() bool
	Run(context.Context) bool
}

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
	// This function returns false if the services cannot be found
	// or the introduction of this edge would produce a cyle.
	Import(ctx context.Context, src, dest id.ServiceId) bool
	// Run blocks the caller until all the prerequistes have been
	// launched.  It returns false if it returned because of
	// a timeout or the service id cannot be found, otherwise true.
	Run(context.Context, id.ServiceId) bool
}

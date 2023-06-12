package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/iansmith/parigot/apishared/id"
	pcontext "github.com/iansmith/parigot/context"
	syscall "github.com/iansmith/parigot/g/syscall/v1"
)

//
// startupService
//

// The default lock discipline for this type is that you should call a method
// on this type when the lock is unlocked.  Any function that
// does assert the lock should make sure to release it before returning.

type startupService struct {
	pkg, name string
	id        id.ServiceId
	runReady  bool
	exported  bool
	started   bool
	parent    *startupCoordinator
	runCh     chan struct{}
	lock      *sync.Mutex
}

// newStartupService creates a representative startupService given a set of
// parameters that define the new service.
func newStartupService(pkg, name string, sid id.ServiceId, parent *startupCoordinator, isClient bool) *startupService {
	result := &startupService{
		pkg:      pkg,
		name:     name,
		id:       sid,
		runReady: false,
		parent:   parent,
		exported: isClient,
		lock:     new(sync.Mutex),
		runCh:    make(chan struct{}),
	}
	_ = Service(result)
	return result
}

// wakeUp causes a send on the servicImpl's runCh and thus checks to see
// if this service can run now.  This method does not lock.
func (s *startupService) wakeUp() {
	s.runCh <- struct{}{}
}

// Id returns the id of this service.
func (s *startupService) Id() id.ServiceId {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.id
}

// Name returns the name, not the fully qualified name, of this service.
func (s *startupService) Name() string {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.name
}

// Package returns the package name, not the fully qualified name, of this service.
func (s *startupService) Package() string {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.pkg
}

// RunRequested returns true if the service has been marked as wanting to run.
// This flag can be true even when the service is not started yet because it can have
// it may have dependencies that are not yet running.
func (s *startupService) RunRequested() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.runReady
}

// RequestRun set this service to be marked as ready to run, and thus is blocked
// waiting until the dependencies are running.
func (s *startupService) RequestRun() {
	s.lock.Lock()
	defer s.lock.Unlock()
	log.Printf("RequestRun: %s%s", s.Name(), s.Short())
	s.runReady = true
}

// Exported returns true if the service has been exported.
func (s *startupService) Exported() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.exported
}

// Run is called to indicate that the service wants to run and can be blocked
// until all of its dependencies are fulfilled.
func (s *startupService) Run(ctx context.Context) syscall.KernelErr {
	s.RequestRun()
	return s.waitToRun(ctx) // be sure this does not lock
}

// canRun is true only if three conditions are met.  The service has been
// exported and has requested to run. The third condition is that the all
// elements in the dependency graph that are "behind" this service are started.
// "Behind" here means that this service may need the service that is behind
// to start running.
func (s *startupService) canRun(ctx context.Context) bool {
	if !s.Exported() {
		return false
	}
	if !s.RunRequested() {
		return false
	}
	pcontext.Debugf(ctx, "trying to see if %s [%s] can run now ", s.Short(), s.Name())
	withFn := pcontext.CallTo(ctx, "checkNodesBehindForRunning")
	depCheck := s.parent.checkNodesBehindForRunning(withFn, s.String())
	return depCheck
}

// waitToRun waits until the timeout expires or until it receives a wake
// up call and a check for the ability to run successfully is made. It returns
// the appropriate KernelErr if there has been a timeout.
//
// waitToRun should be called with the lock available.
func (s *startupService) waitToRun(ctx context.Context) syscall.KernelErr {
	counter := 0
	if s.canRun(ctx) {
		s.SetStarted()
		pcontext.Debugf(ctx, "%s is immediately ready to run", s.Name())
		return syscall.KernelErr_NoError
	}
	pcontext.Debugf(ctx, "Timeout loop running for %s", s.Name())
	for {
		if s.canRun(ctx) {
			s.SetStarted()
			pcontext.Debugf(ctx, "%s is ready to run after waiting (%d)", s.Name(), counter)
			return syscall.KernelErr_NoError
		}
		select {
		case <-s.runCh:
		case <-time.After(1 * time.Second):
			pcontext.Debugf(ctx, "%s incrementing counter: %d", s.Name(), counter)
			counter++
			if counter > timeoutInSecs {
				return syscall.KernelErr_RunTimeout
			}
			// try again
		}
	}
}

// export causes this service to be marked as exported, and this is one
// of the preconditions for starting.
func (s *startupService) export() {
	s.lock.Lock()
	defer s.lock.Unlock()
	log.Printf("export: %s%s", s.Name(), s.Short())
	s.exported = true
}

// String returns a string that is the content of the service id that this startupService
// represents (long form).
func (s *startupService) String() string {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.id.String()
}

// Short returns a string that is the content of the service id that this startupService
// represents (short form).
func (s *startupService) Short() string {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.id.Short()
}

// SetStarted marks a service as running.  This can be done only after this
// startupService has passed the conditions of canRun.
func (s *startupService) SetStarted() {
	s.lock.Lock()
	defer s.lock.Unlock()
	log.Printf("SetStarted: %s%s", s.Name(), s.Short())
	if s.started {
		log.Printf("WARNING! set started called for node that is already started")
	}
	s.started = true
}

// Started returns true if the service has been marked as already running.
func (s *startupService) Started() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.started
}

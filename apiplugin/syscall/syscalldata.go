package main

import "github.com/iansmith/parigot/apishared/id"

var SData = newSyscallDataImpl()

type Service interface {
	Id() id.ServiceId
	Name() string
	Package() string
}

type SyscallData interface {
	//ServiceByName looks up a service and returns it based on the
	//values package_ and name.  If this returns nil, the service could
	//not be found.
	ServiceByName(package_, name string) Service
	// SetService puts a service into SyscallData.  This should only be
	// called once for each package_ and name pair. It returns the
	// ServiceId for the service named, creating a new one if necessary.
	// If the bool result is false, then the pair already existed and
	// we made no changes to it.
	SetService(package_, name string) (Service, bool)
}

type serviceImpl struct {
	pkg, name string
	id        id.ServiceId
}

// Id returns the id of this service.
func (s *serviceImpl) Id() id.ServiceId {
	return s.id
}

// Name returns the name, not the fully qualified name, of this service.
func (s *serviceImpl) Name() string {
	return s.name
}

// Package returns the package name, not the fully qualified name, of this service.
func (s *serviceImpl) Package() string {
	return s.pkg
}

type syscallDataImpl struct {
	sidToService                map[string]Service
	packageNameToServiceNameMap map[string]map[string]Service
	serviceNameToService        map[string]Service
}

func newServiceImpl(pkg, name string, sid id.ServiceId) *serviceImpl {
	result := &serviceImpl{
		pkg:  pkg,
		name: name,
		id:   sid,
	}
	_ = Service(result)
	return result
}
func newSyscallDataImpl() *syscallDataImpl {
	impl := &syscallDataImpl{
		sidToService:                make(map[string]Service),
		packageNameToServiceNameMap: make(map[string]map[string]Service),
		serviceNameToService:        make(map[string]Service),
	}
	_ = SyscallData(impl)
	return impl
}

func (s *syscallDataImpl) SetService(package_, name string) (Service, bool) {
	svc := s.ServiceByName(package_, name)
	if svc != nil {
		return svc, false
	}
	svcId := id.NewServiceId()
	result := newServiceImpl(package_, name, svcId)
	return result, true

}

func (s *syscallDataImpl) ServiceByName(package_, name string) Service {
	nameMap, ok := s.packageNameToServiceNameMap[package_]
	if !ok {
		return nil
	}
	svc, ok := nameMap[name]
	if !ok {
		return nil
	}
	return svc
}

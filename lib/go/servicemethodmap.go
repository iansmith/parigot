package lib

import (
	"context"
	"fmt"

	"github.com/iansmith/parigot/apishared/id"
	"github.com/iansmith/parigot/g/syscall/v1"
	"github.com/iansmith/parigot/lib/go/future"

	"google.golang.org/protobuf/types/known/anypb"
)

// MustRequireFunc is the type of the functions that are created
// by the code generator from protobuf definitions of the
// form MustRequireXXXX(). These are used in the function
// RunXXXX() to indicate required services (dependencies).
type MustRequireFunc func(context.Context, id.ServiceId)

// FuncAnyIO is the type of the guest-side functions that
// implement the set and tear down of method implementations
// in a server.  If a method fleazil is
// defined on a service, it will have FuncAnyIO wrapper
// that unmarshals input parameters and marshals the return
// value.
type FuncAnyIO func(*anypb.Any) future.Method[*anypb.Any, int32]

// Backgrounder is an interface that can be implemented by
// types that want to get period background calls when the
// latest attempt to receive a call has timed out.  Note that
// when the background is actually running (in Background)
// their are no attempts to retreive method calls on this service.
type Backgrounder interface {
	Background(context.Context)
}

// ServiceMethodMap is the data structure that provides conversions
// between a service/method pair and their variants.
// A service or method can be converted to a string with their
// String() method and this method can convert that string back
// to the appropriate service or method. The ServiceMethodMap
// can also convert between the human readable names of
// services and methods and their appropriate Ids. It
// contains a mapping from a service/method pair to a FuncAnyIO
// that is the guest-side implementation of the method.  Finally,
// it allows particular methods to be enabled and disabled so
// they will not be used when generating the list of pairs
// for a call to ReadOne().
type ServiceMethodMap struct {
	forward   map[string]map[string]future.Invoker
	sidString map[string]id.ServiceId
	midString map[string]id.MethodId
	nameToSid map[string]id.ServiceId
	nameToMid map[string]map[string]id.MethodId
	sidToName map[string]string
	midToName map[string]string
	call      []*syscall.ServiceMethodCall
	disabled  map[string]bool
}

func NewServiceMethodMap() *ServiceMethodMap {
	result := &ServiceMethodMap{
		forward:   make(map[string]map[string]future.Invoker),
		sidString: make(map[string]id.ServiceId),
		midString: make(map[string]id.MethodId),
		nameToSid: make(map[string]id.ServiceId),
		nameToMid: make(map[string]map[string]id.MethodId),
		sidToName: make(map[string]string),
		midToName: make(map[string]string),
		//default is false here, so we don't need to add entries
		//unless either it was true before or is true now
		disabled: make(map[string]bool),
	}
	return result
}

const sidMidPairKeyGen = "%s,%s"

// AddServiceMethod is called when a new method has been bound. This
// method creates various data structures needed to be able to look up
// the service and method later, as well as find the appropriate
// FuncAnyIO associated with pair.  Note that the funcAnyIO
// may be nil when the function is not available in this address
// space and any caller must use Dispatch().
func (s *ServiceMethodMap) AddServiceMethod(sid id.ServiceId, mid id.MethodId,
	serviceName, methodName string, fn future.Invoker) {

	methMap, ok := s.forward[sid.String()]
	if !ok {
		s.forward[sid.String()] = make(map[string]future.Invoker)
		methMap = s.forward[sid.String()]
	}
	methMap[mid.String()] = fn

	s.sidString[sid.String()] = sid
	s.midString[mid.String()] = mid

	s.nameToSid[serviceName] = sid
	mMap, ok := s.nameToMid[sid.String()]
	if !ok {
		s.nameToMid[sid.String()] = make(map[string]id.MethodId)
		mMap = s.nameToMid[sid.String()]
	}
	mMap[methodName] = mid

	s.sidToName[sid.String()] = serviceName
	s.midToName[mid.String()] = methodName

	curr := &syscall.ServiceMethodCall{
		ServiceId: sid.Marshal(),
		MethodId:  mid.Marshal(),
	}
	s.call = append(s.call, curr)
}

// Pair returns a list of Service/Method pairs suitable for use in
// a ReadOneRequest.  Particular elements of the map can be omitted
// or included with Disable and Enable.
func (s *ServiceMethodMap) Call() []*syscall.ServiceMethodCall {
	result := []*syscall.ServiceMethodCall{}
	for _, pair := range s.call {
		sidStr := id.UnmarshalServiceId(pair.ServiceId).String()
		midStr := id.UnmarshalMethodId(pair.MethodId).String()
		disableKey := fmt.Sprintf(sidMidPairKeyGen, sidStr, midStr)
		disabled, ok := s.disabled[disableKey]
		if ok && disabled {
			continue
		}
		result = append(result, pair)
	}
	return result
}

// Disable "turns off" a service/method pair within this map.  This
// pair will not appear in results of Pair() until Enable is called for
// this pair.  If the values of this pair of ids does not correspond to
// a real pair that is known to the service, this call is ignored.
// Disable can be useful in production situations where you want disable
// methods on an object that are only for testing.
func (s *ServiceMethodMap) Disable(sid id.ServiceId, mid id.MethodId) {
	disableKey := fmt.Sprintf(sidMidPairKeyGen, sid.String(), mid.String())
	_, ok := s.disabled[disableKey]
	if ok {
		s.disabled[disableKey] = true
	}
}

// Enable "turns on" a given service/method pair within the map.  Thus
// the pair will be returned as part of the Pair() result. If the pair
// of these ids is not found, this call is ignored.
// Enable can be useful in testing situations where you want enable
// methods on an object that are only for testing.
func (s *ServiceMethodMap) Enable(sid id.ServiceId, mid id.MethodId) {
	disableKey := fmt.Sprintf(sidMidPairKeyGen, sid.String(), mid.String())
	_, ok := s.disabled[disableKey]
	if ok {
		delete(s.disabled, disableKey)
	}
}

// Func returns the FuncAnyIO object associated with the sid and mid pair. If
// either sid or mid cannot be found, it returns nil.
func (s *ServiceMethodMap) Func(sid id.ServiceId, mid id.MethodId) future.Invoker {
	m := s.forward[sid.String()]
	if m == nil {
		return nil
	}
	fn, ok := m[mid.String()]
	if !ok {
		return nil
	}
	return fn
}

// MethodNameToId is used to find a method by name, given a particular service id.
// This function returns the value MethodIdZeroValue if either the
// service or the method cannot be found.
func (s *ServiceMethodMap) MethodNameToId(sid id.ServiceId, methodName string) id.MethodId {
	m, ok := s.nameToMid[sid.String()]
	if !ok {
		return id.MethodIdZeroValue()
	}
	mid, ok := m[methodName]
	if !ok {
		return id.MethodIdZeroValue()
	}
	return mid
}

// MethodIdToName is used to find a method given a particular service id.
// This function returns "" if either the service or the method cannot be found.
// This does not require a service id because method ids are unique.
func (s *ServiceMethodMap) MethodIdToName(mid id.MethodId) string {
	name, ok := s.midToName[mid.String()]
	if !ok {
		return ""
	}
	return name
}

// Len returns the number of known methods in this
// ServiceMethodMap
func (s *ServiceMethodMap) Len() int {
	return len(s.call)
}

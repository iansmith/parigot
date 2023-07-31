package wheeler

import (
	"sync"
	"testing"

	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/proto"
)

const package1 = "shaved"
const name1 = "yak"
const package2 = "fleazil"
const name2 = "quux"

func TestExportDirect(t *testing.T) {
	w := newWheeler(nil, nil)
	hid := id.NewHostId()

	//registration first
	sid, ok := registerService(t, w, "some pkg", "some svc", hid)
	if !ok {
		t.Errorf("failed to register service")
		return
	}

	req := makeExportReq(hid, sid)
	a, err := w.export(req)
	if err != syscall.KernelErr_NoError {
		t.Errorf("unexpected kernel error from export: %s",
			syscall.KernelErr_name[int32(err)])
		return
	}
	resp := &syscall.ExportResponse{}
	e := a.UnmarshalTo(resp)
	if e != nil {
		t.Errorf("can't unmarshal to export response: %s", e.Error())
	}

	checkExportStructures(t, w, hid, sid, package1, name1, package2, name2)

}

func TestExportThroughChan(t *testing.T) {

	w := newWheeler(nil, nil)
	hid := id.NewHostId()

	sid, ok := registerService(t, w, "bleah", "yech", hid)
	if !ok {
		t.Errorf("unable to register service")
		return
	}

	req := makeExportReq(hid, sid)
	ch := makeRequstToWheeler(w, req)
	result := <-ch

	resp := &syscall.ExportResponse{}
	err := result.A.UnmarshalTo(resp)
	if err != nil {
		t.Errorf("unable to unmarshal to export response: %s", err.Error())
	}
	if result.Err != syscall.KernelErr_NoError {
		t.Errorf("unexpected error in response to export response: %s", syscall.KernelErr_name[int32(result.Err)])
	}
	checkExportStructures(t, w, hid, sid, package1, name1, package2, name2)

}

func TestExitThroughChan(t *testing.T) {
	exitCh := make(chan int32)

	wg := sync.WaitGroup{}
	wg.Add(1)
	// this is required because the sender (wheeler) will not
	//
	monitorFn := func(c chan int32) {
		// no block because it should be already sent
		x := <-c
		if x != 74 {
			t.Errorf("unexpected value read from exit channel: %d", x)
		}
		wg.Add(-1)

	}
	go monitorFn(exitCh)

	w := newWheeler(nil, exitCh)
	req := &syscall.ExitRequest{Code: 74}

	ch := makeRequstToWheeler(w, req)
	result := <-ch

	resp := &syscall.ExitResponse{}
	err := result.A.UnmarshalTo(resp)

	if err != nil {
		t.Errorf("unexpected  error from UnmarshalTo: %s", err.Error())
	}

	if result.Err != syscall.KernelErr_NoError {
		t.Errorf("unexpected  error from Exit call: %s", syscall.KernelErr_name[int32(result.Err)])
	}
	// should return pretty much immediately because the reader and
	// writer have already run
	wg.Wait()
}

func TestBindMethodThroughChan(t *testing.T) {
	w := newWheeler(nil, nil)
	req := &syscall.BindMethodRequest{}
	resp := &syscall.BindMethodResponse{}

	hid := id.NewHostId()

	sid, ok := registerService(t, w, "baz", "quux", hid)
	if !ok {
		t.Errorf("unable to register service")
		return
	}
	expReq := makeExportReq(hid, sid)
	if _, err := w.export(expReq); err != syscall.KernelErr_NoError {
		t.Errorf("unable to create service for bind test: %s", syscall.KernelErr_name[int32(err)])
		return
	}

	req.HostId = hid.Marshal()
	req.ServiceId = sid.Marshal()
	req.MethodName = "fleazil"

	if makeReqRespToWheeler(t, w, req, resp) != syscall.KernelErr_NoError {
		return
	}

	mid := id.UnmarshalMethodId(resp.GetMethodId())
	checkServiceMethodBinding(t, w, sid, req.MethodName, mid)
}

func TestRegisterThroughChannel(t *testing.T) {
	w := newWheeler(nil, nil)
	req := &syscall.RegisterRequest{}
	resp := &syscall.RegisterResponse{}

	hid := id.NewHostId()
	pkg := "bleah"
	name := "blah"
	req.Fqs = &syscall.FullyQualifiedService{
		PackagePath: pkg,
		Service:     name,
	}
	req.HostId = hid.Marshal()

	if kerr := makeReqRespToWheeler(t, w, req, resp); kerr != syscall.KernelErr_NoError {
		t.Errorf("unable to successfully make req to wheeler (%s)", syscall.KernelErr_name[int32(kerr)])
		return
	}
	if resp.ExistedPreviously {
		t.Errorf("expected existed previously to be false")
	}
	newSid := id.UnmarshalServiceId(resp.Id)
	checkRegistrationStructures(t, w, pkg, name, newSid, hid)

	// no effect expect in the existed previously field
	if makeReqRespToWheeler(t, w, req, resp) != syscall.KernelErr_NoError {
		return
	}
	if !resp.ExistedPreviously {
		t.Errorf("expected existed previously to be true")
	}
	checkRegistrationStructures(t, w, pkg, name, newSid, hid)
}

func TestRegistrationAndLookup(t *testing.T) {

	w := newWheeler(nil, nil)
	req := &syscall.RegisterRequest{}
	resp := &syscall.RegisterResponse{}

	pkg := "foo"
	name := "Foo"
	hid := id.NewHostId()

	expectNoService(t, w, pkg, name)

	msg := "expected service to not exist yet: %s.%s"
	req.Fqs = &syscall.FullyQualifiedService{
		PackagePath: pkg,
		Service:     name,
	}
	req.HostId = hid.Marshal()

	if err := makeReqRespToWheeler(t, w, req, resp); err != syscall.KernelErr_NoError {
		t.Errorf("unable to register %s.%s", pkg, name)
		return
	}
	svc := id.UnmarshalServiceId(resp.GetId())
	existed := resp.ExistedPreviously
	if existed {
		t.Errorf(msg, pkg, name)
	}
	msg = "returned service id is either zero or empty valuet: %s.%s"
	if svc.IsZeroOrEmptyValue() {
		t.Errorf(msg, pkg, name)
	}

	msg = "expected service to already exist: %s.%s"
	if err := makeReqRespToWheeler(t, w, req, resp); err != syscall.KernelErr_NoError {
		t.Errorf("unable to make request to register with wheeler: %s,%s", pkg, name)
		return
	}
	if !resp.ExistedPreviously {
		t.Errorf(msg, pkg, name)
	}

	expectConsistentResult(t, w, svc, pkg, name, hid)
}

func TestStartupSimple(t *testing.T) {
	//none of these would actually work, this is to test
	// the startup logic
	w := newWheeler(nil, nil)

	aHost := id.NewHostId()
	bHost := id.NewHostId()
	cHost := id.NewHostId()
	dHost := id.NewHostId()
	host := []id.HostId{aHost, bHost, cHost, dHost}
	letter := []string{"a", "b", "c", "d"}
	svc := make([]id.ServiceId, len(letter))

	// a exports x.x
	export := [][]*syscall.FullyQualifiedService{
		{
			{
				PackagePath: "x",
				Service:     "x",
			},
		},
		// b exports x.x and y.y
		{
			{
				PackagePath: "x",
				Service:     "x",
			},
			{
				PackagePath: "y",
				Service:     "y",
			},
		},
		// c exports z.z
		{
			{
				PackagePath: "z",
				Service:     "z",
			},
		},
		// d exports "z.z" and "w.w"
		{
			{
				PackagePath: "z",
				Service:     "z",
			},
			{
				PackagePath: "w",
				Service:     "w",
			},
		},
	}

	// topo sort is
	// C, D, B, A

	// a requires y.y and w.w
	require := [][]*syscall.FullyQualifiedService{
		{
			{
				PackagePath: "y",
				Service:     "y",
			},
			{
				PackagePath: "w",
				Service:     "w",
			},
		},
		// b requires w.w and z.z
		{
			{
				PackagePath: "w",
				Service:     "w",
			},
			{
				PackagePath: "z",
				Service:     "z",
			},
		},
		// c requires nothing
		nil,
		// d requires z.z
		{
			{
				PackagePath: "z",
				Service:     "z",
			},
		},
	}

	for i, l := range letter {
		regReq := &syscall.RegisterRequest{
			Fqs: &syscall.FullyQualifiedService{
				PackagePath: l,
				Service:     l,
			},
			HostId: host[i].Marshal(),
		}
		regResp := &syscall.RegisterResponse{}
		if err := makeReqRespToWheeler(t, w, regReq, regResp); err != syscall.KernelErr_NoError {
			t.Errorf("unable to make registration request for %s:%s", l, syscall.KernelErr_name[int32(err)])
			return
		}
		svc[i] = id.UnmarshalServiceId(regResp.GetId())
	}
	for i := 0; i < len(letter); i++ {
		exReq := &syscall.ExportRequest{
			ServiceId: svc[i].Marshal(),
			HostId:    host[i].Marshal(),
		}
		exReq.Service = export[i]

		exResp := &syscall.ExportResponse{}
		if err := makeReqRespToWheeler(t, w, exReq, exResp); err != syscall.KernelErr_NoError {
			t.Errorf("unable to make export request: %s", syscall.KernelErr_name[int32(err)])
			return
		}
	}
	for i := 0; i < len(letter); i++ {
		reqReq := &syscall.RequireRequest{
			Source: svc[i].Marshal(),
			Dest:   require[i],
		}
		reqResp := &syscall.RequireResponse{}
		if err := makeReqRespToWheeler(t, w, reqReq, reqResp); err != syscall.KernelErr_NoError {
			t.Errorf("unable to make export request: %s", syscall.KernelErr_name[int32(err)])
			return
		}
	}

}

//
// HELPERS
//

func checkExportStructures(t *testing.T, w *wheeler, hid id.HostId, sid id.ServiceId,
	pkg1, name1, pkg2, name2 string) {

	checkHostToService(t, w, hid, sid)
	checkServiceToHost(t, w, hid, sid)

	checkPkgToService(t, w, pkg1, name1, sid, hid)
	checkPkgToService(t, w, pkg2, name2, sid, hid)

}

func checkPkgToService(t *testing.T, w *wheeler, pkg, name string, sid id.ServiceId, hid id.HostId) {
	t.Helper()

	s, ok := w.pkgToServiceImpl[pkg]
	if !ok {
		t.Errorf("unable to find package %s", pkg)
		return
	}
	bind, ok := s[name]
	if !ok {
		t.Errorf("unable to find named service %s inside package %s", pkg, name)
		return
	}
	if len(bind) != 1 {
		t.Errorf("unable to find correct bindings for %s.%s, length is %d but should be 1",
			pkg, name, len(bind))
		return
	}
	if !bind[0].host.Equal(hid) {
		t.Errorf("bad host id found as binding for %s.%s, expected %s but got %s",
			pkg, name, hid.String(), bind[0].host.String())
		return
	}
	if !bind[0].service.Equal(sid) {
		t.Errorf("bad service id found as binding for %s.%s, expected %s but got %s",
			pkg, name, sid.String(), bind[0].service.String())
		return
	}
}

func checkHostToService(t *testing.T, w *wheeler, hid id.HostId, sid id.ServiceId) {
	t.Helper()
	// validate construction of internal data structs
	ok := true
	var svc []id.ServiceId
	_, ok = w.hostToService[sid.String()]
	if ok {
		t.Errorf("should not have the service id as a key in the hostToService")
	}
	svc, ok = w.hostToService[hid.String()]
	if !ok {
		t.Errorf("did not find host key in the hostToService")
	}
	if len(svc) != 1 {
		t.Errorf("wrong length for svcs of host, expected 1 but got %d", len(svc))
		return
	}
	if !svc[0].Equal(sid) {
		t.Errorf("wrong service id found in hostToService, expected %s but got %s",
			sid.String(), svc[0].String())
	}
}

func checkServiceToHost(t *testing.T, w *wheeler, hid id.HostId, sid id.ServiceId) {
	t.Helper()
	cand, ok := w.serviceToHost[sid.String()]
	if !ok {
		t.Errorf("unable to find sid %s in service to host", sid.String())
		return
	}
	if !cand.Equal(hid) {
		t.Errorf("wrong host for service %s, found %s but expected %s",
			sid.String(), hid.String(), cand.String())
	}
}

func makeExportReq(hid id.HostId, sid id.ServiceId) *syscall.ExportRequest {
	return &syscall.ExportRequest{
		ServiceId: sid.Marshal(),
		HostId:    hid.Marshal(),
		Service: []*syscall.FullyQualifiedService{
			{
				PackagePath: package1,
				Service:     name1,
			},
			{
				PackagePath: package2,
				Service:     name2,
			},
		},
	}
}

func makeRequstToWheeler(w *wheeler, msg proto.Message) chan OutProtoPair {
	toWh := w.ch
	ch := make(chan OutProtoPair)
	inPair := InProtoPair{
		Msg: msg,
		Ch:  ch,
	}
	toWh <- inPair
	return ch
}
func makeReqRespToWheeler(t *testing.T, w *wheeler, req proto.Message, resp proto.Message) syscall.KernelErr {
	ch := makeRequstToWheeler(w, req)
	result := <-ch

	if result.Err != syscall.KernelErr_NoError {
		return result.Err
	}
	e := result.A.UnmarshalTo(resp)
	if e != nil {
		t.Errorf("unable to unmarshal result of bind method:%s", e.Error())
		return syscall.KernelErr_MarshalFailed
	}
	return syscall.KernelErr_NoError
}

func checkServiceMethodBinding(t *testing.T, w *wheeler, sid id.ServiceId, name string, mid id.MethodId) {
	t.Helper()

	methMap, ok := w.serviceToMethMap[sid.String()]
	if !ok {
		t.Errorf("unable to find methods for service %s", sid.String())
		return
	}
	cand, ok := methMap[name]
	if !ok {
		t.Errorf("unable to find method %s for service %s", name, sid.String())
		return
	}
	if !cand.Equal(mid) {
		t.Errorf("mismatch of method ids, expected %s but got %s", mid.String(), cand.String())
		return
	}

}

func checkRegistrationStructures(t *testing.T, w *wheeler, pkg, name string, sid id.ServiceId, hid id.HostId) {
	nameMap, ok := w.pkgToRegistration[pkg]
	if !ok {
		t.Errorf("unable to find package %s in registration map", pkg)
		return
	}
	bind, ok := nameMap[name]
	if !ok {
		t.Errorf("unable to find service %s.%s in registration map", pkg, name)
		return
	}
	if len(bind) != 1 {
		t.Errorf("wrong number of host bindings for registration of  %s.%s, expected 1 but got %d", pkg, name, len(bind))
		return
	}
	if !bind[0].host.Equal(hid) {
		t.Errorf("bad host id for registration of service %s.%s, expected %s but got %s", pkg, name,
			hid.String(), bind[0].host.String())
	}
	if !bind[0].service.Equal(sid) {
		t.Errorf("bad host id for registration of service %s.%s, expected %s but got %s", pkg, name,
			sid.String(), bind[0].service.String())
	}
}

func expectNoService(t *testing.T, w *wheeler, pkg, name string) {
	t.Helper()
	msg := "found service when not expecting it, %s.%s"
	req := &syscall.ServiceByNameRequest{}
	resp := &syscall.ServiceByNameResponse{}
	req.Fqs = &syscall.FullyQualifiedService{
		PackagePath: pkg,
		Service:     name,
	}

	if kerr := makeReqRespToWheeler(t, w, req, resp); kerr != syscall.KernelErr_NoError {
		t.Errorf("unable to make request to ServiceByName in wheeler: %s", syscall.KernelErr_name[int32(kerr)])
		return
	}
	msg = "should never match an unknown service: %s.%s"
	if len(resp.Binding) != 0 {
		t.Errorf(msg, pkg, name)
	}

	sid := id.ServiceIdZeroValue()
	msg = "should never match an zero service id: %s"
	reqId := &syscall.ServiceByIdRequest{}
	respId := &syscall.ServiceByIdResponse{}
	reqId.ServiceId = sid.String()
	if kerr := makeReqRespToWheeler(t, w, reqId, respId); kerr != syscall.KernelErr_NoError {
		t.Errorf("unable to make request to ServiceById in wheeler: %s", syscall.KernelErr_name[int32(kerr)])
		return
	}
	msg = "should never match an empty service id: %s"
	if len(resp.Binding) != 0 {
		t.Errorf(msg, sid.String())
	}

	sid = id.NewServiceId()
	reqId.ServiceId = sid.String()
	msg = "should not match a new, random service id: %s"
	if kerr := makeReqRespToWheeler(t, w, reqId, respId); kerr != syscall.KernelErr_NoError {
		t.Errorf("unable to make request to ServiceById in wheeler: %s", syscall.KernelErr_name[int32(kerr)])
		return
	}
	if len(resp.Binding) != 0 {
		t.Errorf(msg, sid.String())
	}
}
func expectConsistentResult(t *testing.T, w *wheeler, sid id.ServiceId, pkg, name string, hid id.HostId) {
	t.Helper()

	req := &syscall.RegisterRequest{
		Fqs: &syscall.FullyQualifiedService{
			PackagePath: pkg,
			Service:     name,
		},
	}
	req.HostId = hid.Marshal()
	resp := &syscall.RegisterResponse{}
	if err := makeReqRespToWheeler(t, w, req, resp); err != syscall.KernelErr_NoError {
		t.Errorf("unable to make register request for %s.%s (%s)", pkg, name, syscall.KernelErr_name[int32(err)])
		return
	}
	svc2 := id.UnmarshalServiceId(resp.GetId())

	msg := "expected service id to be consistent in two uses of register: %s.%s"
	if !sid.Equal(svc2) {
		t.Errorf(msg, pkg, name)
	}

	msg = "expected service id to be consistent betweenf SetService and ServiceById: %s.%s"
	reqById := &syscall.ServiceByIdRequest{
		ServiceId: sid.String(),
	}
	respById := &syscall.ServiceByIdResponse{}
	if err := makeReqRespToWheeler(t, w, reqById, respById); err != syscall.KernelErr_NoError {
		t.Errorf("unable to find service by id for %s", sid.String())
		return
	}
	svc3 := id.UnmarshalServiceId(resp.GetId())
	if !svc3.Equal(sid) {
		t.Errorf(msg, pkg, name)
		return
	}
	msg = "expected service id to be consistent betweenf RegisterService and ServiceByName: %s.%s"
	reqByName := &syscall.ServiceByNameRequest{
		Fqs: &syscall.FullyQualifiedService{
			PackagePath: pkg,
			Service:     name,
		},
	}
	respByName := &syscall.ServiceByNameResponse{}
	if err := makeReqRespToWheeler(t, w, reqByName, respByName); err != syscall.KernelErr_NoError {
		t.Errorf("unable to find service by name for %s.%s", pkg, name)
		return
	}
	if len(respByName.GetBinding()) != 1 {
		t.Errorf("wrong number of results for service by name %s.%s", pkg, name)
		return
	}
	svc4 := id.UnmarshalServiceId(respByName.GetBinding()[0].ServiceId)
	host4 := id.UnmarshalHostId(respByName.GetBinding()[0].HostId)
	if !svc4.Equal(sid) {
		t.Errorf(msg, pkg, name)
		return
	}

	if !host4.Equal(hid) {
		t.Errorf("failed to match expected host id when looking up service by name %s.%s", pkg, name)
		return
	}
}

func registerService(t *testing.T, w *wheeler, pkg, name string, hid id.HostId) (id.ServiceId, bool) {
	regReq := &syscall.RegisterRequest{
		Fqs: &syscall.FullyQualifiedService{
			PackagePath: pkg,
			Service:     name,
		},
	}
	regReq.HostId = hid.Marshal()
	regResp := &syscall.RegisterResponse{}
	a, err := w.register(regReq)
	if err != syscall.KernelErr_NoError {
		t.Errorf("uable to call register directly: %s", syscall.KernelErr_name[int32(err)])
		return id.ServiceIdZeroValue(), false
	}
	if err := a.UnmarshalTo(regResp); err != nil {
		t.Errorf("unable to unmarshal register response : %s", err.Error())
		return id.ServiceIdZeroValue(), false
	}
	sid := id.UnmarshalServiceId(regResp.GetId())
	return sid, true

}

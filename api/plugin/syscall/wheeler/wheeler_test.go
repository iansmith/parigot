package wheeler

import (
	"testing"

	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/g/syscall/v1"
)

const package1 = "shaved"
const name1 = "yak"
const package2 = "fleazil"
const name2 = "quux"

func TestExportDirect(t *testing.T) {

	w := newWheeler(nil)
	sid := id.NewServiceId()
	hid := id.NewHostId()

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

	w := newWheeler(nil)
	sid := id.NewServiceId()
	hid := id.NewHostId()

	toWh := w.ch
	req := makeExportReq(hid, sid)
	ch := make(chan OutProtoPair)
	inPair := InProtoPair{
		Msg: req,
		Ch:  ch,
	}
	toWh <- inPair
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

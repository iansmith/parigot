package future

import (
	"testing"

	"github.com/iansmith/parigot/g/syscall/v1"
)

func TestSimpleBaseWithValue(t *testing.T) {
	value := uint64(27)
	f := NewBaseWithValue[uint64](value)
	var output uint64

	// since the value is already resolved in the future we
	// should get this call immediately
	f.Handle(func(u uint64) {
		output = u
	})

	if value != output {
		t.Errorf("%s expected %d but got %d", "TestSimpleBaseWithValue", value, output)
	}
}

func TestSimpleBaseHandle(t *testing.T) {
	value := int16(52)
	var output int16
	f := NewBase[int16]()
	f.Handle(func(x int16) {
		output = x
	})
	// no change in output yet
	if output != 0 {
		t.Errorf("%s (0): expected %d but got %d", "TestSimpleBaseHandle", 0, output)
	}
	// set the value in the promise
	f.Set(value)
	// now check the output
	if output != value {
		t.Errorf("%s (1): expected %d but got %d", "TestSimpleBaseHnadle", value, output)
	}
}

func TestBaseMultipleHandle(t *testing.T) {
	value := 2
	var out1, out2, out3 int
	future := NewBase[int]()
	future.Handle(func(v int) {
		out1 = v + 2
	})
	future.Handle(func(v int) {
		out2 = v * 10
	})
	ranAny := future.Set(value)
	if !ranAny {
		t.Errorf("%s (0): future.Set should return true if handle functions are run", "TestBaseMultipleHandle")
	}
	if out1 != value+2 {
		t.Errorf("%s (1): expected %d but got %d", "TestBaseMultipleHandle", value+2, out1)
	}
	if out2 != value*10 {
		t.Errorf("%s (2): expected %d but got %d", "TestBaseMultipleHandle", value*10, out2)
	}
	out2 = 902
	// since the others have already run, this will run immediately
	future.Handle(func(v int) {
		out3 = v - 12
	})
	if out3 != value-12 {
		t.Errorf("%s (3):expected %d but got %d", "TestBaseMultipleHandle", value-12, out3)
	}
	// no more handlers to run, so set should have no effect
	if future.Set(10000) != false {
		t.Errorf("%s (5):expected false when no handle functions run", "TestBaseMultipleHandle")
	}
	// make sure no damage to existing values
	if (out1 != value+2) ||
		(out2 != 902) ||
		(out3 != value-12) {
		t.Errorf("%s (7): expected all three values to be unchanged when an 'extra' set is called", "TestBaseMultipleHandle")
	}
}

func TestHandleLater(t *testing.T) {
	value := 16
	var out int

	future := NewBase[int]()
	future.Set(value)
	if out != 0 {
		t.Errorf("%s (0):expected no change to variable", "TestHandleLater")
	}
	future.HandleLater(func(v int) {
		out = v * 2
	})
	if out != 0 {
		t.Errorf("%s (1):expected no change to variable", "TestHandleLater")
	}
	future.Set(2)
	if out != 4 {
		t.Errorf("%s (2):expected %d but got %d", "TestHandleLater", 4, out)
	}
}

func TestMethodSimple(t *testing.T) {
	success := false
	var result *syscall.DispatchResponse
	fail := syscall.KernelErr(0)
	sample := &syscall.DispatchResponse{}

	s := func(d *syscall.DispatchResponse) {
		success = true
		result = d
	}
	f := func(kerr syscall.KernelErr) {
		fail = kerr
	}
	m := NewMethod(s, f)
	//no change before the call
	if success || fail != 0 {
		t.Errorf("%s (0):expected no change before CompleteMethod()", "TestMethodSimple")
	}
	m.CompleteMethod(nil, syscall.KernelErr_DependencyCycle)
	if success || fail != syscall.KernelErr_DependencyCycle {
		t.Errorf("%s (1):expected %s but got %s for error", "TestMethodSimple", syscall.KernelErr_name[int32(syscall.KernelErr_DependencyCycle)],
			syscall.KernelErr_name[int32(fail)])
	}
	//
	// try success case
	//
	success = false
	fail = 0
	m = NewMethod(s, f)
	m.CompleteMethod(sample, syscall.KernelErr_NoError)
	if !success || result != sample || fail != 0 {
		t.Errorf("%s (2):expected success to be true and the sample value to be propagated (%v,%v)", "TestMethodSimple", success, result == sample)
	}
}

func TestMethodReplaceFunc(t *testing.T) {
	success := false
	err := syscall.KernelErr_NoError
	sample := &syscall.DispatchResponse{}

	s := func(d *syscall.DispatchResponse) {
		success = true
	}
	m := NewMethod(func(*syscall.DispatchResponse) {},
		func(_ syscall.KernelErr) {})
	// s should replace no op we started with
	m.Success(s)
	if success {
		t.Errorf("%s (0):premature call to success", "TestMethodReplaceFunc")

	}
	m.CompleteMethod(sample, syscall.KernelErr_NoError)
	if !success || err != 0 {
		t.Errorf("%s (1):expected replacement func to set success", "TestMethodReplaceFunc")
	}
	// try failure case
	success = false
	m = NewMethod(func(*syscall.DispatchResponse) {},
		func(_ syscall.KernelErr) {})
	f := func(kerr syscall.KernelErr) {
		err = kerr
	}
	m.Failure(f)
	if success || err != 0 {
		t.Errorf("%s (2):call to failure too soon", "TestMethodReplaceFunc")
	}
	m.CompleteMethod(nil, syscall.KernelErr_ReadOneTimeout)
	if success || err != syscall.KernelErr_ReadOneTimeout {
		t.Errorf("%s (3):call to change err value", "TestMethodReplaceFunc")
	}
}

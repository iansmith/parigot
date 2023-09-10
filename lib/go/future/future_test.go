//go:build wasip1

package future

import (
	"context"
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
	//test that we don't mess with them too soon
	if success || fail != 0 {
		t.Errorf("%s (0):expected no change before CompleteMethod()", "TestMethodSimple")
	}
	m.CompleteMethod(context.Background(), nil, syscall.KernelErr_DependencyCycle)
	if success || fail != syscall.KernelErr_DependencyCycle {
		t.Errorf("%s (1):expected %s but got %s for error (%v,%d)", "TestMethodSimple", syscall.KernelErr_name[int32(syscall.KernelErr_DependencyCycle)],
			syscall.KernelErr_name[int32(fail)], success, fail)
	}
	//
	// try success case
	//
	success = false
	fail = 0
	m = NewMethod(s, f)
	m.CompleteMethod(context.Background(), sample, syscall.KernelErr_NoError)
	if !success || result != sample || fail != 0 {
		t.Errorf("%s (2):expected success to be true and the sample value to be propagated (%v,%v)", "TestMethodSimple", success, result == sample)
	}
}

func TestMethodQueueFunc(t *testing.T) {
	success1, success2 := false, false
	result := &syscall.DispatchResponse{}
	err1, err2 := syscall.KernelErr_NoError, syscall.KernelErr_NoError
	sample := &syscall.DispatchResponse{}

	m := NewMethod[*syscall.DispatchResponse, syscall.KernelErr](func(resp *syscall.DispatchResponse) { success1, result = true, nil },
		func(kerr syscall.KernelErr) { err1 = kerr })
	m.Success(func(resp *syscall.DispatchResponse) {
		result = resp // note that this will overwrite prev value
		success2 = true
	})
	// the two success should be stacked (really queued)
	m.CompleteMethod(context.Background(), sample, syscall.KernelErr_NoError)
	if !success1 || !success2 || err1 != 0 || err2 != 0 || result != sample {
		t.Errorf("%s (1):expected queued funcs to set both success vars and the result var, as well as not touch errors", "TestMethodQueueFunc")
	}

	// try failure case
	success2, success1 = false, false
	result = nil
	m = NewMethod(func(*syscall.DispatchResponse) { success1 = true },
		func(kerr syscall.KernelErr) { err1 = kerr })

	f := func(kerr syscall.KernelErr) {
		err2 = kerr
	}
	m.Failure(f) // stacks both error funcs
	m.CompleteMethod(context.Background(), sample, syscall.KernelErr_BadId)
	if success1 || success2 || result != nil {
		t.Errorf("%s (2):call to failure should only mess with the two error values", "TestMethodQueueFunc")
	}
	if err1 != syscall.KernelErr_BadId || err2 != syscall.KernelErr_BadId {
		t.Errorf("%s (3):error values should be Bad Id (%d) but are err1=%d, err2=%d", "TestMethodQueueFunc",
			syscall.KernelErr_BadId, err1, err2)
	}
}

func TestCancelMethod(t *testing.T) {
	success1, success2 := false, false
	err1 := syscall.KernelErr_NoError
	result := &syscall.DispatchResponse{}
	sample := &syscall.DispatchResponse{}

	m := NewMethod[*syscall.DispatchResponse, syscall.KernelErr](func(resp *syscall.DispatchResponse) { success1, result = true, nil },
		func(kerr syscall.KernelErr) { err1 = kerr })
	m.Success(func(resp *syscall.DispatchResponse) {
		result = resp // note that this will overwrite prev value
		success2 = true
	})
	m.Cancel()
	m.CompleteMethod(context.Background(), sample, 0)
	if success1 || success2 || result == sample || err1 != 0 {
		t.Errorf("%s (0): all queued calls to Success should have removed by Cancel", "TestCancelMethod")
	}
	m.CompleteMethod(context.Background(), nil, 1)
	if success1 || success2 || result == sample || err1 != 0 {
		t.Errorf("%s (1): all queued calls to Failure should have removed by Cancel", "TestCancelMethod")
	}
}

func TestCancelBase(t *testing.T) {
	counter := 0

	base := NewBase[uint8]()
	base.Handle(func(_ uint8) {
		counter++
	})
	base.Handle(func(_ uint8) {
		counter++
	})
	base.Handle(func(_ uint8) {
		counter++
	})
	base.Cancel()
	base.Set(0xff)
	if counter != 0 {
		t.Errorf("%s (0): all queued calls to Success should have removed by Cancel", "TestCancelBase")
	}
	base.Handle(func(_ uint8) {
		counter = 1000
	})
	if counter != 1000 {
		t.Errorf("%s (1): expected counter to be 1000 but is %d, expected it to change Handle() on already completed Base should run the func immediately", "TestCancelBase", counter)
	}
}
func TestAllSuccess(t *testing.T) {
	successX3 := false
	allSuccess := false
	sample := &syscall.DispatchRequest{}
	var ptr *syscall.DispatchRequest
	var x1, x2, x3 *Method[*syscall.DispatchRequest, syscall.KernelErr]

	x1, x2, x3, fut := setupAll(&successX3, &ptr)
	fut.Success(func() {
		allSuccess = true
	})
	x1.CompleteMethod(context.Background(), sample, 0)
	x2.CompleteMethod(context.Background(), sample, 0)
	x3.CompleteMethod(context.Background(), sample, 0)
	//now the fut should be finished too
	if !successX3 || !allSuccess || ptr != sample {
		t.Errorf("%s (0): expected completing all the update success markers and ptr", "TestAllSuccess")
	}
}
func TestAllFail(t *testing.T) {
	sample := &syscall.DispatchRequest{}
	successX3 := false
	var ptr *syscall.DispatchRequest
	var x1, x2, x3 *Method[*syscall.DispatchRequest, syscall.KernelErr]
	allFail := false
	failedIndex := -1

	x1, x2, x3, fut := setupAll(&successX3, &ptr)
	fut.Success(func() {
		t.Errorf("%s (0): should not call success when x3 fails", "TestAllFail")
	})
	fut.Failure(func(index int) {
		allFail = true
		failedIndex = index
	})
	x1.CompleteMethod(context.Background(), sample, 0)
	x2.CompleteMethod(context.Background(), sample, 0)
	x3.CompleteMethod(context.Background(), nil, syscall.KernelErr_BadId)
	if !allFail || failedIndex != 2 || successX3 {
		t.Errorf("%s (1): x3 failed, expected changed values for failure only", "TestAllFail")
	}

}

// attempt to refactor may have made code above worse
func setupAll(successX3 *bool, ptr **syscall.DispatchRequest) (*Method[*syscall.DispatchRequest, syscall.KernelErr],
	*Method[*syscall.DispatchRequest, syscall.KernelErr], *Method[*syscall.DispatchRequest, syscall.KernelErr], *AllFuture[*syscall.DispatchRequest, syscall.KernelErr]) {
	*successX3 = false
	*ptr = nil

	// make sure it this doesn't get called
	s := func(d *syscall.DispatchRequest) {
		*successX3 = true
		*ptr = d
	}

	x1 := NewMethod[*syscall.DispatchRequest, syscall.KernelErr](nil, nil)
	x2 := NewMethod[*syscall.DispatchRequest, syscall.KernelErr](nil, nil)
	x3 := NewMethod[*syscall.DispatchRequest, syscall.KernelErr](s, nil)

	fut := All(x1, x2, x3)
	return x1, x2, x3, fut
}

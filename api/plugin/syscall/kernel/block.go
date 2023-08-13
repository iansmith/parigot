package kernel

import (
	"math"
	"reflect"
	"time"

	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/g/syscall/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// ReadOne computes a vector of different things to listen to.  It will listen
// for the minimum of the maxWait variable and the smallest time given by any of
// the components.  The copmonents are the receiver and the listener channels,
// represented by the Receiver and Finisher types.
func (k *kdata) ReadOne(req *syscall.ReadOneRequest, resp *syscall.ReadOneResponse) syscall.KernelErr {
	unified := make([]reflect.SelectCase, maxCases)
	count := 0
	minTimeoutMillis := math.MaxInt

	//top priority is launches
	lcb, remaining := k.start.Ready()
	if !lcb.sid.IsZeroValue() {
		// we mess with the matcher() here to complete the launch
		// request (and fire its future) which will be discovered by
		// the code later that asks for ready from the matcher
		launchErr := syscall.KernelErr_NoError
		if lcb.hasCycle {
			launchErr = syscall.KernelErr_DependencyCycle
		}
		a := &anypb.Any{}
		if err := a.MarshalFrom(&syscall.LaunchResponse{}); err != nil {
			klog.Errorf("unable to marshal launch response: %v", err)
		} else {
			k.matcher().Response(lcb.cid, a, int32(launchErr))
		}
		klog.Infof("launch completed for %s, %d remaining", lcb.sid.Short(), remaining)
	}

	// we favor completing futures over reading in new requests
	// that may be a terrible idea
	hid := id.UnmarshalHostId(req.HostId)
	err := k.responseReady(hid, resp)
	if err != syscall.KernelErr_NoError {
		return err
	}
	if resp.Resolved != nil {
		// we got a resolution and we are done
		return syscall.KernelErr_NoError
	}

	// we want to check all the receivers
	for _, r := range k.rawRecv {
		c := r.Ch()
		if c != nil {
			selectCase := reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(c)}
			unified[count] = selectCase
			count++
		}
		if r.TimeoutInMillis() < minTimeoutMillis {
			if r.TimeoutInMillis() < 0 {
				klog.Errorf("ignoring bad timeout value: %d", r.TimeoutInMillis())
			} else {
				minTimeoutMillis = r.TimeoutInMillis()
			}
		}
		if count == maxCases {
			klog.Errorf("maximum number of possible network receive ports reached (%d), ignoring excess", maxCases)
			break
		}
	}

	// we normally add 2 special cases
	if count < maxCases-2 {
		// setup the cancel chan
		unified[count] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(k.cancel)}
		count++
		// setup the timeout
		if minTimeoutMillis == math.MaxInt {
			minTimeoutMillis = maxWait
		}
		timeoutCh := time.After(time.Duration(minTimeoutMillis) * time.Millisecond)
		unified[count] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(timeoutCh)}
		count++
	} else {
		klog.Errorf("unable to add special cases for timeout or cancel")
	}

	// the actual select
	chosen, value, ok := reflect.Select(unified[:count])
	if !ok {
		klog.Errorf("unexpected close of channel in ReadOne(): %T", unified[chosen])
		return syscall.KernelErr_ChannelClosed
	}
	// timeout?
	if chosen == len(unified)-1 {
		return syscall.KernelErr_ReadOneTimeout
	}
	// cancel?
	if chosen == len(unified)-2 {
		// somebody trying to stop this method running
		return syscall.KernelErr_NoError
	}

	//
	// All readers work the same way
	//
	selectReturn := value.Interface().(*syscall.ReadOneResponse)
	proto.Merge(resp, selectReturn)
	return syscall.KernelErr_NoError
}
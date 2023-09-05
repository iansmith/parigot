package lib

import (
	"context"
	"time"

	syscallguest "github.com/iansmith/parigot/api/guest/syscall"
)

var parigotTime = "parigot_time"
var parigotTimeIndex = "parigot_time_index"

// CurrentTime returns the current time in the current timezone unless the current
// time has been set with SetCurrentTime. If that is the case, it returns the
// values provided to SetCurrentTime, in order.  Note that the returned context.Context
// should be used as a replacement for the context provided as a parameter,
// because the state of the context changed.
//
// All guest code should use this call rather than time.Now() or similar
// to provide a convenient way to test time-dependent behavior.
func CurrentTime(ctx context.Context) (context.Context, time.Time) {
	raw := ctx.Value(parigotTime)
	if raw != nil {
		t, ok := raw.([]time.Time)
		if ok {
			i := ctx.Value(parigotTimeIndex).(int)
			next := (i + 1) % (len(t))
			result := context.WithValue(ctx, parigotTimeIndex, next)
			now := t[i]
			return result, now
		}
	}
	// we got no value, just use clock
	return ctx, time.Now().In(syscallguest.CurrentTimezone())
}

var timeIndex int
var timeSeq []time.Time

// SetCurrentTime should only be used in tests. Calling this method _freezes_
// time so that future calls will see a chosen time (in t) as the current
// time.  Every call to CurrentTime() receives the next value in the sequence
// provided by t.  If there are more calls than values in t, the sequence is
// restarted from the beginning, thus providing a single value "stops time"
// at that point.
//
// The returned context.Context should be used as a replacement for the context
// provided as a parameter.
func SetCurrentTime(ctx context.Context, t ...time.Time) context.Context {
	if len(t) == 0 {
		panic("no time values provided to SetCurrentTime")
	}
	return context.WithValue(context.WithValue(ctx, parigotTime, t),
		parigotTimeIndex, 0)
}

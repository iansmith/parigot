package lib

import logmsg "github.com/iansmith/parigot/g/msg/log/v1"

// This interfae is kind of a hack around an import cycle.  Lib can't import proto/g/log because
// that results in a sequence of imports that will lead to problems with syscall.  Underlying
// prob is that everybody uses lib to get access Id handling.

type Log interface {
	// AbortOnFatal() bool
	// SetAbortOnFatal(bool)
	Log(req *logmsg.LogRequest) error
}

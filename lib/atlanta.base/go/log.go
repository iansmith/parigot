package lib

import (
	pblog "github.com/iansmith/parigot/api/proto/g/pb/log"
)

// This interfae is kind of a hack around an import cycle.  Lib can't import proto/g/log because
// that results in a sequence of imports that will lead to problems with syscall.  Underlying
// prob is that everybody uses lib to get access Id handling.

type Log interface {
	// AbortOnFatal() bool
	// SetAbortOnFatal(bool)
	Log(prefix string, level pblog.LogLevel, msg string)
}

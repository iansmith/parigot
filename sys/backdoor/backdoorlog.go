// This entire package is here to avoid import loops.  Many different parts of
// the kernel and service implementations need to use the backdoor to the
// logger. However, that implementation is part of a larger logging implementation
// and imports of that lead to frequent import loops.
//
// This package is here to allow any part of the _go_ implementation of kernel
// or services to have a way to import and use this interface without needing
// to worry about import loops.  The implementation of this interface is
// passed into this packages SetInternalLogger() by an init() function of
// the logging package (github.com/iansmith/parigot/api_impl/log/go_).
package backdoor

import (
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
)

// InternalLogger just specifies the log interface that is visible to
// the code that _implements_ this functionality.
type InternalLogger interface {
	ProcessLogRequest(*logmsg.LogRequest, bool, bool, bool, []byte)
}

// logger is initialized  by the init() function over in the logging package.
var logger InternalLogger

// SetInternalLogger is called  by the init() function over in the logging package.
func SetInternalLogger(il InternalLogger) {
	logger = il
}

// Log should *only* be called by parts of the implementation of parigot that need
// logging facilities. The isKernel and isBackend should be set to true only if this is called by some part of the kernel itself
// or some "backend" implementation of a service, respectively.  When the isJS flag is set, it indicates
// that the original caller was in wasm and did `log.Printf()` or some other use of log (which is discouraged
// but we have to handle it due to existing go code in the library).  The isJS flag is proper subset of the
// isBackend since literally a backend service implementation calls this, but the isJS flag lets the decoration
// of the log output be more clear.  If the caller does not have the already serialized
// version of req, buffer can be passed as nil and this function will create the buffer itself.
func Log(req *logmsg.LogRequest, isKernel, isBackend bool, isJS bool, buffer []byte) {
	logger.ProcessLogRequest(req, isKernel, isBackend, isJS, buffer)
}

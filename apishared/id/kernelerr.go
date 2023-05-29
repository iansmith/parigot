package id

const (
	// KernelServiceNamespaceExhausted is returned when the kernel can no
	// along accept additional packages, services, or methods.  This is used
	// primarily to thwart attempts at DOS attacks.
	KernelNamespaceExhausted KernelErrIdCode = iota + KernelErrIdGuestStart
	// KernelNotFound means that a package, service, or method that was requested
	// could not be found.
	KernelNotFound
	// KernelDataTooLarge means that the size of some part of remote call was bigger
	// than the buffer allocated to receive it.  This could be a problem either on the call or
	// the return.
	KernelDataTooLarge
	// KernelMarshalFailed is an internal error of the kernel. This means that
	// a marshal of a protobuf has failed.  This is only used in situations
	// that are internel to the kernel--if user code misbehaves in this fashion
	// an error is sent to the program _from_ the kernel about the failure.
	KernelMarshalFailed
	// KernelCallerUnavailable means that the kernel could not find the original caller
	// requeted the computation for which results have been provided.  It is most likely
	// because the caller was killed, exited or timed out.
	KernelCallerUnavailable
	// KernelServiceAlreadyClosedOrExported means that some process has already reported
	// the service in question as closed or has already expressed that it is
	// exporting (implementing this service).  This is very likely a case where there
	// are two servers that think they are or should be implementing the same service.
	KernelServiceAlreadyClosedOrExported
	// KernelServiceAlreadyRequired means that this same process has already
	// required the given service.
	KernelServiceAlreadyRequired
	// KernelDependencyCycle means that no deterministic startup ordering
	// exists for the set of exports and requires in use.  In other words,
	// you must refactor your program so that you do not have a cyle to make
	// it come up cleanly.
	KernelDependencyCycle
	// KernelNetworkFailed means that we successfully connected to the nameserver, but failed
	// during the communication process itself.
	KernelNetworkFailed
	// KernelNetworkConnectionLost means that our internal connection to the remote nameserver
	// was either still working but has lost "sync" in the protocol or the connection has
	// become entirely broken.  The kernel will close the connection to remote nameserver
	// and reestablish it after this error.
	KernelNetworkConnectionLost
	// KernelDataTooSmall means that the kernel was speaking some protocol with a remote server,
	// such as a remote nameserver, and data read from the remote said was smaller than the protocol
	// dictated, e.g. it did not contain a checksum after a data block.
	KernelDataTooSmall
	// KernelConnectionFailed means that the attempt to open a connection to a remote
	// service has failed to connect.
	KernelConnectionFailed
	// KernelNSRetryFailed means that we tried twice to reach the nameserver with
	// the given request, but both times could not do so.
	KernelNSRetryFailed
	// KernelDecodeError indicates that an attempt to extract a protobuf object
	// from an encoded set of bytes has failed.  Typically, this means that
	// the encoder was not called.
	KernelDecodeError
	// KernelExecError means that we received a response from the implenter of a particular
	// service's function and the execution of that function failed.
	KernelExecError
	// KernelBadId means received something from your code that was supposed to be an error and
	// it did not have the proper mark on it (IsErrorType()).
	KernelBadId
	// KernelDependencyFailure means that the dependency infrastructure has failed.  This is different
	// than when a user creates bad set of depedencies (KernelDependencyCycle).
	KernelDependencyFailure
	// KernelAbortRequest indicates that the program that receives this error
	// should exit because the nameserver has asked it to do so.  This
	// means that some _other_ program has failed to start correctly, so this
	// deployment cannot succeed.
	KernelAbortRequest
	// KernelExitRequest indicates that the program that receives this error
	// should exit because the nameserver has asked it to do so.  This is not really
	// an "error" but rather an indication that the program that requested the
	// exit may do so immediately.
	KernelExitRequest
	// KernelEncodeError indicates that an attempt encode a protobuf
	// with header and CRC has failed.
	KernelEncodeError
	// KernelClosedErr indicates that that object is now closed.  This is used
	// as a signal when writing data between the guest and host.
	KernelClosedErr
	// KernelGuestReadFailed indicates that we did not successfully read
	// from guest memory. This is usually caused by the address read at being
	// out of bounds.
	KernelGuestReadFailed
	// KernelGuestWriteFailed indicates that we did not successfully write
	// to guest memory. This is usually caused by the address written at being
	// out of bounds.
	KernelGuestWriteFailed
)

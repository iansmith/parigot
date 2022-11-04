package lib

// GetMaxMessageSize is the maximum size of a buffer that is allowed
// to result from marshalling parameters or return results.  Callers
// of remote functions should allocate this amount of space when they
// are getting a result back.
func GetMaxMessageSize() int32 {
	return 0x400 // 4k
}

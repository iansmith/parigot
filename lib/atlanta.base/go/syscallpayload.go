package lib

// RegPayload is what the WASM code sends to the kernel when Register()
// is called.  It's the same content as LocatePayload right now, since one
// is write (RegPayload) and the other is read.
type RegPayload struct {
	PkgPtr       int64     // in p0a
	PkgLen       int64     // in p0b
	ServicePtr   int64     // in p1a
	ServiceLen   int64     // in p1b
	ErrorPtr     *[2]int64 // out p0
	ServiceIdPtr *[2]int64 // out p1
}

// LocatePayload is what the WASM code sends to the kernel when Locate()
// is called.  It's the same content as RegPayload right now, since one
// is write and the other is read (LocatePayload).
type LocatePayload struct {
	PkgPtr       int64     // in p0a
	PkgLen       int64     // in p0b
	ServicePtr   int64     // in p1a
	ServiceLen   int64     // in p1b
	ErrorPtr     *[2]int64 // out p0
	ServiceIdPtr *[2]int64 // out p1
}

// DispatchPayload is what the WASM code sends to the kernel when Dispatch()
// is called.
type DispatchPayload struct {
	ServiceId [2]int64 // in p0
	MethodPtr int64    // in p1a
	MethodLen int64    // in p1b
	CallerPtr int64    // in p2a
	CallerLen int64    // in p2b
	PctxPtr   int64    // in p3a
	PctxLen   int64    //in p3b -- we don't send the capacity
	ParamPtr  int64    // in p4a
	ParamLen  int64    //in p4b -- we don't send the capacity

	// we know the type of this pointer,so we can unmarshal it
	// no problem.  The following one is the problem. We need it's length
	// because all we can do is use the any as a proxy for it.
	OutPctxPtr int64     // out p0 -- serialized protobuf
	ResultPtr  int64     // out p1 -- serialized protobuf
	ErrorPtr   *[2]int64 // out p2

	// ResultLen is both in *and* out parameter. Going in it tells the
	// kernel how much space we allocated for the result. If the call
	// is successful (and the result will fit) it then holds the actual
	// size of the result.
	ResultLen int64 // in p5, out p3
	// OutPctxLen is both in *and* out parameter. Going in it tells the
	// kernel how much space we allocated for the result. If the call
	// is successful (and the result will fit) it then holds the actual
	// size of the result.
	OutPctxLen int64 // in p6, out p4
}

// BindPayload is sent to the kernel to register a particular function
// as the implementation of the tuple (pkg,service,method).
type BindPayload struct {
	PkgPtr     int64     // in p0a
	PkgLen     int64     // in p0b
	ServicePtr int64     // in p1a
	ServiceLen int64     // in p1b
	MethodPtr  int64     // in p2a
	MethodLen  int64     // in p2a
	Direction  int64     // in p3
	MethodId   *[2]int64 // out p0
	ErrorPtr   *[2]int64 // out p1
}

// BlockPayload is sent to the kernel to request the server block the caller
// until a method request for the server is ready.  The results of this
// function are placed in the buffers provided.
type BlockPayload struct {
	PctxPtr   int64     // in p0a
	PctxLen   int64     // in p0b <-- also out p0
	ParamPtr  int64     // in p1a
	ParamLen  int64     // in p1b <--- also out p1
	Direction int64     // in p2 ///xxxxshould go away?
	MethodId  *[2]int64 // out p2
	CallId    *[2]int64 // out p3
	ErrorPtr  *[2]int64 // out p4
}

// ReturnValuePayload is sent to the kernel request the server send these results
// back to the caller as the results of their previous call.  Note that the CallId
// allows the original caller to verify this is the right response.
type ReturnValuePayload struct {
	PctxPtr        int64    // in p0a
	PctxLen        int64    // in p0b
	ResultPtr      int64    // in p1a
	ResultLen      int64    // in p1b
	MethodId       [2]int64 // in p2
	CallId         [2]int64 // in p3
	KernelErrorPtr [2]int64 // in p4 <--- also out p0

	//xxx should be doing this too
	//UserErrorPtr   *[2]int64 // out p5
}

// Export payload is used to indicate to the server what
// services you implement.  No more method implementations
// can be bound once this method is called, so you want to
// be sure you've finished your BindMethod() calls before
// you use this.  This should only be called by a server.
type ExportPayload struct {
	PkgPtr         int64     // in p0a
	PkgLen         int64     // in p0b
	ServicePtr     int64     // in p1a
	ServiceLen     int64     // in p1b
	KernelErrorPtr *[2]int64 // out p0

}

// RequirePayload is use to indicate what services you consume.
// This function may be called by servers, but is always called
// by clients.
type RequirePayload ExportPayload

// Start payload is used by clients and servers to indicate that
// they are ready to start.  This call blocks until the server has
// start all the prerequisites (the services requested via RequirePayload)
// and thus there is no need to wait to use Locate().  This call
// can fail if there is a dependency cycle.  If it does it will
// try to write the loop problem into the Loop result buffer here
// but will not do so if the loop printout is larger than LoopResultLen.
// Each segment of the loop is separated by ; in the result buffer.
type StartPayload struct {
	KernelErrorPtr *[2]int64 // out p0
	LoopResultPtr  int64     // out p1
	LoopResultLen  int64     // in p0 and out p2
}

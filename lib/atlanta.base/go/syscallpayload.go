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
	PkgPtr     int64 // in p0a
	PkgLen     int64 // in p0b
	ServicePtr int64 // in p1a
	ServiceLen int64 // in p1b
	MethodPtr  int64 // in p2a
	MethodLen  int64 // in p2a
	FuncPtr    int64 // in p3
}

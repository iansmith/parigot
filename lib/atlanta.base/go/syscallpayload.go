package lib

// RegPayload is what the WASM code sends to the kernel when Register()
// is called.  It's the same content as LocatePayload right now, since one
// is write (RegPayload) and the other is read.
type RegPayload struct {
	PkgPtr          int64     // in p0a
	PkgLen          int64     // in p0b
	ServicePtr      int64     // in p1a
	ServiceLen      int64     // in p1b
	OutErrPtr       *[2]int64 // out p0
	OutServiceIdPtr *[2]int64 // out p1
}

// LocatePayload is what the WASM code sends to the kernel when Locate()
// is called.  It's the same content as RegPayload right now, since one
// is write and the other is read (LocatePayload).
type LocatePayload struct {
	PkgPtr          int64     // in p0a
	PkgLen          int64     // in p0b
	ServicePtr      int64     // in p1a
	ServiceLen      int64     // in p1b
	OutErrPtr       *[2]int64 // out p0
	OutServiceIdPtr *[2]int64 // out p1
}

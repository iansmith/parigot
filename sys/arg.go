package sys

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

// The (golang+wasm) 	linker guarantees global data starts from at least wasmMinDataAddr.
// Keep in sync with cmd/link/internal/ld/data.go:wasmMinDataAddr.
const wasmStartAddr = 4096
const wasmMinDataAddr = wasmStartAddr + 8192

type ptrPair struct {
	val string // just for debugging help
	ptr int32  // pointer into the wasm memory
}

func updateBuffer(buffer *bytes.Buffer, s string) int32 {
	asByte := []byte(s)
	buffer.Write(asByte)
	buffer.Write([]byte{0})

	result := (int32(len(asByte)) + 1)
	if result%8 != 0 {
		for i := int32(0); i < 8-(result%8); i++ {
			buffer.Write([]byte{0})
		}
		result += 8 - (result % 8)
	}
	return result
}

// https://github.com/golang/go/blob/db36eca33c389871b132ffb1a84fd534a349e8d8/misc/wasm/wasm_exec.js#L488
func GetBufferFromArgsAndEnv(m Service, startOfArgs int32) (*bytes.Buffer, int32, error) {
	entries := len(m.GetArg()) + len(m.GetEnv()) + 3
	all := make([]ptrPair, entries) // two null strings and argv[0]
	buffer := &bytes.Buffer{}

	count := 0
	offset := int32(startOfArgs)
	for _, a := range append([]string{m.GetName()}, m.GetArg()...) {
		all[count].ptr = offset
		all[count].val = a
		offset += updateBuffer(buffer, a)
		count++
	}
	all[count].ptr = 0
	count++
	for _, envvar := range m.GetEnv() {
		all[count].ptr = offset
		all[count].val = envvar
		offset += updateBuffer(buffer, envvar)
		count++
	}
	all[count].ptr = 0
	count++

	buf := make([]byte, 8)
	argv := offset

	for _, pair := range all[:count] {
		//log.Printf("xxx send all %d: %x (%s)", i, pair.ptr, pair.val)
		buf[4], buf[5], buf[6], buf[7] = 0, 0, 0, 0
		binary.LittleEndian.PutUint32(buf[:4], uint32(pair.ptr))
		//log.Printf("buffer bytes: %x", buf)
		buffer.Write(buf)
	}
	log.Printf("final buffer len %d", buffer.Len())
	if int32(buffer.Len())+startOfArgs >= wasmMinDataAddr {
		return nil, 0, fmt.Errorf("microservice %s has args+environment size of %d bytes, but max is %d",
			m.GetName(), buffer.Len()-wasmStartAddr, wasmMinDataAddr-wasmStartAddr)
	}
	return buffer, argv, nil
}

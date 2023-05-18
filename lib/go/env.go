package lib

import (
	"bytes"
	"encoding/binary"
	"log"
	"strings"
	"unsafe"
)

// All of this file's code is a workaround.  We load the linear memory with what seems like it should be
// correctly formatted argv and envp.  We also call the run() function with two args to start up a process,
// but the args and the environment are nowhere to be found. We lay out the memory the same way the
// JS-hosted code does, so it appears there is some code that is not run during our initialization that
// is run by the JS side.
//
// We do one horrible hack though and we store argv (the char**) at a fixed place in memory (0x1000) and
// pull that value to bootstrap the process of pulling out the argv and envp vectors.  This means that
// if you want "byte for byte" similarly with the JS-hosted code, you'd have to drop that 4 byte value.
//

// loadCString returns a go string from a sequence of bytes that is a C-style (null terminated) string.
// This used, effectively, to convert a C char* to a go string.
func loadCString(rawAddr int32) string {
	var buf bytes.Buffer
	i := int32(0)
	for {
		str := (*byte)(unsafe.Pointer(uintptr(int32(rawAddr + i))))
		buf.WriteByte(*str)
		if *str == 0 {
			break
		}
		i++
	}
	asBytes := buf.Bytes()
	return string(asBytes[:len(asBytes)-1])

}

// loadArgvPointer is a hack. We store the argv value as pointer at address 0x1000.
func loadArgvPointer(argv int32) int32 {
	var buf [4]byte
	for i := int32(0); i < 4; i++ {
		str := (*byte)(unsafe.Pointer(uintptr(int32(argv + i))))
		log.Printf("xxxx loadArgvPointer: %d, %p", i, str)
		buf[i] = *str
	}
	return int32(binary.LittleEndian.Uint32(buf[:]))
}

// computeArgvEnvp reads the _sequence_ of pointers (each 4 bytes wide) that starts at argv and
// then uses each of these to read a C style string.  The args and the env are terminated by an empty
// string.
func computeArgvEnvp() ([]string, []string) {
	log.Printf("compute ArgvEnvp")
	// xxx hacky: we put argv at a fixed memory location because we have no way to pass param to this function... except via argv!
	argvAddr := loadArgvPointer(0x1000 - 4)
	log.Printf("compute ArgvEnvp2")
	argVector := []string{}
	envVector := []string{}
	index := int32(0)
	log.Printf("xxx computeArgvEnvp333333")
	for {
		ptr := loadArgvPointer(argvAddr + (8 * index))
		index++
		if ptr == 0 {
			break
		}
		argVector = append(argVector, loadCString(ptr))
	}
	for {
		ptr := loadArgvPointer(argvAddr + (8 * index))
		index++
		if ptr == 0 {
			break
		}
		envVector = append(envVector, loadCString(ptr))
		index++
	}
	log.Printf("xxx computeArgvEnvy2")

	return argVector, envVector
}

var envp, argv []string

func FlagParseCreateEnv() {
	// log.Printf("xxx -- FlagParseCreateEnv()")
	// argv, envp = computeArgvEnvp()
	// os.Args = argv
	// flag.Parse()
}

// This is a workalike for os.Getenv()
func Getenv(envvar string) string {
	for _, candidate := range envp {
		part := strings.Split(candidate, "=")
		if len(part) != 2 {
			// this is checked for at the time the config is read
			panic("badly formed environment, variable " + envvar)
		}
		if part[0] == envvar {
			return part[1]
		}
	}
	return ""
}

// This is a workalike for os.LookupEnv().  It can be used to differentiate an empty, but set,
// environment variable from the an enviroment variable that is simply not present.
func LookupEnv(envvar string) (string, bool) {
	for _, candidate := range envp {
		part := strings.Split(candidate, "=")
		if len(part) != 2 {
			// this is checked for at the time the config is read
			panic("badly formed environment, variable " + envvar)
		}
		if part[0] == envvar {
			return part[1], true
		}
	}
	return "", false
}

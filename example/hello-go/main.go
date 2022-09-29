package main

import (
	"github.com/iansmith/parigot/lib/base/go/log"
)

//export main.main
func main() {
	var l = log.LocalT{}
	l.SetLogMask(log.DevMask)
	l.SetAbortOnFatal(true)

	l.Debug("hello %s", "go")
}

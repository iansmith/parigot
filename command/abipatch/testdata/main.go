package main

import (
	"github.com/iansmith/parigot/command/jsstrip/testdata/ex1"
	"github.com/iansmith/parigot/lib/base/go/log"
)

func main() {
	log.Dev.Debug("hello, logger")
	ex1.Driver()
	log.Dev.Debug("goodbye, logger")
}

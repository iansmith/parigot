package main

import (
	"github.com/iansmith/parigot/g/parigot/abi"
	"github.com/iansmith/parigot/lib/base/go/log"
)

func main() {
	log.Dev.Debug("hello, logger")
	abi.Exit(1)
}

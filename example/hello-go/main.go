package main

import (
	"github.com/iansmith/parigot/abi/go/abi"
	"github.com/iansmith/parigot/lib/base/go/log"
)

func main() {
	log.Dev.Debug("hello, logger")
	abi.Exit(1)
}

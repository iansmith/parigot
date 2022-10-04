package main

import (
	"github.com/iansmith/parigot/abi/go/abi"
)

func main() {
	abi.OutputString("bleah")
	abi.Exit(107)
}

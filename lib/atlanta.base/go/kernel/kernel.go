// This package is a thin wrapper around kernel functionality so you can call
// that functionality easily and without having to create a KernelClient each time.
// It allows you to use the notation kernel.Exit(0), for example, if you import this
// under the default name.
package kernel

import (
	"github.com/iansmith/parigot/g/parigot/kernel"
)

var kernelClient *kernel.KernelClient

func initSvc() {
	var err error
	kernelClient, err = kernel.LocateKernel()
}

func Exit(code int32) {
	if kernelClient==nil {
		initSvc()
	}
	kernelClient.
}

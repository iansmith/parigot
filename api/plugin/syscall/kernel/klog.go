package kernel

import (
	"log"
	"os"
)

var klog KLog = newKernelLogger()

type kernelLogger struct {
	*log.Logger
}

func (k *kernelLogger) Errorf(spec string, rest ...interface{}) {
	k.Logger.Printf("ERR : "+spec, rest...)
}
func (k *kernelLogger) Warnf(spec string, rest ...interface{}) {
	k.Logger.Printf("WARN: "+spec, rest...)
}
func (k *kernelLogger) Infof(spec string, rest ...interface{}) {
	k.Logger.Printf("INFO: "+spec, rest...)
}

func newKernelLogger() KLog {
	return &kernelLogger{
		Logger: log.New(os.Stdout, "kernel:", log.Default().Flags()),
	}
}

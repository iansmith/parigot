package kernel

import (
	"fmt"
	"log/slog"
	"os"
)

var klog KLog = newKernelLogger()

type kernelLogger struct {
	*slog.Logger
}

func (k *kernelLogger) Errorf(spec string, rest ...interface{}) {
	k.Logger.Error(fmt.Sprintf(spec, rest...))
}
func (k *kernelLogger) Warnf(spec string, rest ...interface{}) {
	k.Logger.Warn(fmt.Sprintf(spec, rest...))
}
func (k *kernelLogger) Infof(spec string, rest ...interface{}) {
	k.Logger.Info(fmt.Sprintf(spec, rest...))
}

func newKernelLogger() KLog {
	opt := &slog.HandlerOptions{}
	opt.Level = slog.LevelWarn
	l := slog.New(slog.NewTextHandler(os.Stdout, opt)).With("kernel", true)
	return &kernelLogger{Logger: l}
}

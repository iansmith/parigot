package kernel

import (
	"log/slog"
	"os"
)

var klog KLog = newKernelLogger()

type kernelLogger struct {
	*slog.Logger
}

func (k *kernelLogger) Errorf(spec string, rest ...interface{}) {
	k.Logger.Error(spec, rest...)
}
func (k *kernelLogger) Warnf(spec string, rest ...interface{}) {
	k.Logger.Warn(spec, rest)
}
func (k *kernelLogger) Infof(spec string, rest ...interface{}) {
	k.Logger.Info(spec, rest...)
}

func newKernelLogger() KLog {
	opt := &slog.HandlerOptions{}
	opt.Level = slog.LevelWarn
	l := slog.New(slog.NewTextHandler(os.Stdout, opt)).With("kernel", true)
	return &kernelLogger{Logger: l}
}

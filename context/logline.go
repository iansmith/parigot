package context

import (
	"context"

	"github.com/fatih/color"
)

type logLine struct {
	source Source
	level  LogLevel
	data   [MaxLineLen]byte // c-style terminator (nul byte)
}

var defaultColor *color.Color

func init() {
	defaultColor = color.New(color.FgBlack)
	if !UseBlack {
		defaultColor = color.New(color.FgWhite)
	}
}

func LogLineFromString(ctx context.Context, line string, src Source, lvl LogLevel) *logLine {
	if src == UnknownS {
		if v := ctx.Value(ParigotSource); v != nil {
			src = v.(Source)
		}
	}
	maxWithZero := MaxLineLen - 1
	if len(line) >= maxWithZero {
		start := len(line) - maxWithZero
		line = line[:start]
	}
	i := 0
	result := &logLine{}
	for i < len(line) {
		result.data[i] = line[i]
		i++
	}
	result.level = lvl
	result.source = src
	return result
}

func (l *logLine) color() func(a ...interface{}) string {
	if l.level == stackTraceInternal {
		return defaultColor.Add(color.Faint).SprintFunc()
	}
	switch l.source {
	case Source(UnknownS):
		return addLogLevelVisual(color.New(color.BgGreen), l.level).SprintFunc()
	case Source(Client):
		return addLogLevelVisual(color.New(color.FgGreen), l.level).SprintFunc()
	case Source(ServerGo):
		return addLogLevelVisual(color.New(color.FgYellow), l.level).SprintFunc()
	case Source(ServerWasm):
		return addLogLevelVisual(color.New(color.FgHiYellow), l.level).SprintFunc()
	case Source(Parigot):
		return addLogLevelVisual(color.New(color.FgHiCyan), l.level).SprintFunc()
	}
	panic("unable to understand source of logline")
}
func addLogLevelVisual(c *color.Color, l LogLevel) *color.Color {
	switch l {
	case Fatal:
		return c.Add(color.BlinkSlow)
	case Error:
		return c.Add(color.ReverseVideo)
	case Warn:
		return c.Add(color.Underline)
	case Debug:
		return c.Add(color.Faint)
	}
	return c
}

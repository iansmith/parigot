package context

import (
	"context"
	"fmt"
	"sync"

	"github.com/fatih/color"
)

type logLine struct {
	source         Source
	level          LogLevel
	funcName, spec string
	value          []interface{}
	data           [MaxLineLen]byte // c-style terminator (nul byte)
	raw            bool
	lock           *sync.Mutex
	prevCtx        context.Context
}

var defaultColor *color.Color
var oppDefaultColor *color.Color

var maxStrLenWithoutColor = 240 // 256 - 16

func init() {
	defaultColor = color.New(color.FgHiBlack)
	oppDefaultColor = color.New(color.FgHiWhite)
	if !UseBlack {
		defaultColor = color.New(color.FgHiWhite)
		oppDefaultColor = color.New(color.FgHiBlack)

	}
}

func NewLogLine(ctx context.Context, src Source, lvl LogLevel, funcName string,
	raw bool, spec string, rest ...interface{}) *logLine {
	src = PullSource(ctx, UnknownS)
	result := &logLine{}
	result.level = lvl
	result.source = src
	result.funcName = pullFunc(ctx, funcName)
	result.raw = raw
	result.spec = spec
	result.lock = new(sync.Mutex)
	result.value = rest
	result.prevCtx = ctx
	return result
}

func (ll *logLine) Print(ctx context.Context) {
	ll.lock.Lock()
	defer ll.lock.Unlock()

	var line string
	if ll.raw {
		line += ll.spec
	} else {
		prefix := detailPrefix(ll.prevCtx, ll.level, ll.source, ll.funcName)
		if ll.spec == "" {
			line = fmt.Sprintf("%s", ll.value[0])
		} else {
			if len(ll.spec) == 0 {
				line = "\n"
			} else {
				line = fmt.Sprintf(ll.spec, ll.value...)
				line += "\n"
			}
			if len(line) > maxStrLenWithoutColor {
				diff := len(line) - maxStrLenWithoutColor
				line = line[diff:]
			}
		}
		line = prefix + line
	}
	var baseColor *color.Color
	switch ll.source {
	case Source(UnknownS):
		baseColor = oppDefaultColor
	case Source(Client):
		baseColor = color.New(color.BgGreen)
	case Source(ServerGo):
		baseColor = color.New(color.FgYellow)
	case Source(ServerWasm):
		baseColor = color.New(color.FgHiYellow)
	case Source(Parigot):
		baseColor = color.New(color.FgCyan)
	case Source(Wazero):
		baseColor = color.New(color.FgBlue)
	case Source(WasiOut):
		baseColor = defaultColor
	case Source(WasiErr):
		baseColor = color.New(color.FgRed)
	}
	mod := addLogLevelVisual(baseColor, ll.level)
	str := mod.SprintfFunc()(line)
	fmt.Print(str)
	//mod.Print(line)
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

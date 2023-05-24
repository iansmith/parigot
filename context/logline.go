package context

import (
	"context"
	"fmt"

	"github.com/fatih/color"
)

type logLine struct {
	source         Source
	level          LogLevel
	funcName, spec string
	value          []interface{}
	data           [MaxLineLen]byte // c-style terminator (nul byte)
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

func LogLineFromPrintf(ctx context.Context, src Source, lvl LogLevel, funcName, spec string, rest ...interface{}) *logLine {
	if src == UnknownS {
		if v := ctx.Value(ParigotSource); v != nil {
			src = v.(Source)
		}
	}
	result := &logLine{}
	result.level = lvl
	result.source = src
	result.funcName = funcName
	result.spec = spec
	result.value = rest
	return result
}

func (ll *logLine) Print() {
	formatted := fmt.Sprintf(ll.spec, ll.value...)
	if len(formatted) == 0 {
		formatted = "\n"
	} else {
		if len(ll.spec) > 1 && ll.spec[len(ll.spec)-1] != '\n' {
			if formatted[len(formatted)-1] != '\n' {
				formatted = formatted + "\n"
			}
		}
	}
	if len(formatted) > maxStrLenWithoutColor {
		diff := len(formatted) - maxStrLenWithoutColor
		formatted = formatted[diff:]
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
	baseColor.Print(formatted)
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

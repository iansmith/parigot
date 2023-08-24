package guest

import (
	"bytes"
	"context"
	"log/slog"

	"github.com/iansmith/parigot/api/shared/id"

	"github.com/fatih/color"
)

const LoggerCtxKey = "logger_context_key"

var colorSeq []color.Attribute = []color.Attribute{
	color.FgHiGreen,
	color.FgHiYellow,
	color.FgHiBlue,
	color.FgHiMagenta,
	color.FgHiCyan,
	color.FgGreen,
	color.FgYellow,
	color.FgBlue,
	color.FgMagenta,
	color.FgCyan,
}

var handlerCount = 0

type ParigotHandler struct {
	h   slog.Handler
	buf *bytes.Buffer
	clr *color.Color
}

func NewParigotHandler(sid id.ServiceId) slog.Handler {
	buf := &bytes.Buffer{}

	th := slog.NewTextHandler(buf, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	})
	h := th.WithAttrs([]slog.Attr{
		slog.String("service", sid.Short()),
	})

	result := &ParigotHandler{
		h:   h,
		clr: color.New(colorSeq[handlerCount%len(colorSeq)]),
		buf: buf,
	}
	handlerCount++
	return result
}

func (p *ParigotHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (p *ParigotHandler) Handle(ctx context.Context, rec slog.Record) error {
	p.buf.Reset()
	if err := p.h.Handle(ctx, rec); err != nil {
		return err
	}
	p.clr.Print(p.buf.String())
	return nil
}

func (p *ParigotHandler) WithAttrs(attr []slog.Attr) slog.Handler {
	p.h = p.h.WithAttrs(attr)
	return p
}

func (p *ParigotHandler) WithGroup(g string) slog.Handler {
	p.h = p.h.WithGroup(g)
	return p
}

package guest

import (
	"context"
	"log/slog"
	"os"

	"github.com/iansmith/parigot/api/shared/id"
)

const LoggerCtxKey = "logger_context_key"

var handlerCount = 0

type ParigotHandler struct {
	h slog.Handler
}

func NewParigotHandler(sid id.ServiceId) slog.Handler {

	th := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})
	h := th.WithAttrs([]slog.Attr{
		slog.String("service", sid.Short()),
	})

	result := &ParigotHandler{
		h: h,
	}
	handlerCount++
	return result
}

func (p *ParigotHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (p *ParigotHandler) Handle(ctx context.Context, rec slog.Record) error {
	return p.h.Handle(ctx, rec)
}

func (p *ParigotHandler) WithAttrs(attr []slog.Attr) slog.Handler {
	p.h = p.h.WithAttrs(attr)
	return p
}

func (p *ParigotHandler) WithGroup(g string) slog.Handler {
	p.h = p.h.WithGroup(g)
	return p
}

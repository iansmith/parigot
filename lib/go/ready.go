package lib

import (
	"context"

	"github.com/iansmith/parigot/api/shared/id"
	"github.com/iansmith/parigot/lib/go/future"
)

type ReadyChecker interface {
	Ready(ctx context.Context, sid id.ServiceId) *future.Base[bool]
}

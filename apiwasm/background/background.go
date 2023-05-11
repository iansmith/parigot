package background

import "context"

type Background interface {
	Background(ctx context.Context)
}

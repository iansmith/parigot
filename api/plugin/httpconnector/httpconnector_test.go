package httpconnector

import (
	"context"
	"testing"

	pcontext "github.com/iansmith/parigot/context"
	"github.com/iansmith/parigot/g/httpconnector/v1"
)

func TestCheck(t *testing.T) {
	svc := newHttpCntSvc(context.Background())
	req := &httpconnector.CheckRequest{
		Method: "test",
	}
	resp := &httpconnector.CheckResponse{}
	ctx := pcontext.DevNullContext(context.Background())

	svc.check(ctx, req, resp)
}

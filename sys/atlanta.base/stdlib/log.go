package stdlib

import (
	"log"

	"github.com/iansmith/parigot/g/pb/kernel"
	pblog "github.com/iansmith/parigot/g/pb/log"
	"github.com/iansmith/parigot/g/pb/parigot"
	"github.com/iansmith/parigot/lib"
)

func init() {
	impl := newLogImpl()
	req := &kernel.BindMethodRequest{
		ProtoPackage: "", //xxx fixme, should be named in parigot package
		Service:      "log",
		Method:       "Log",
	}
	_, err := lib.BindMethod(req, impl.Log)
	if err != nil {
		log.Fatalf("unable to bind log service: %v", err)
	}

}

type logImpl struct {
}

func newLogImpl() *logImpl {
	return &logImpl{}
}

func (l *logImpl) Log(pctx parigot.PCtx, message *pblog.LogRequest) {
	log.Printf("xxxx log called! %+v, %+v", pctx, message)
}

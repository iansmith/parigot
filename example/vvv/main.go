package main

import (
	"demo/vvv/proto/g/vvv/pb"

	"github.com/iansmith/parigot/g/pb/kernel"
	"github.com/iansmith/parigot/g/pb/parigot"
	"github.com/iansmith/parigot/lib"
	"google.golang.org/protobuf/proto"
)

type StoreServer interface {
	BestOfAllTime(pctx *parigot.PCtx) (proto.Message, error)
	Revenue(pctx *parigot.PCtx, in proto.Message) (proto.Message, error)
	SoldItem(pctx *parigot.PCtx, in proto.Message) error
}

func main() {
	for {
		methodBuf := make([]byte, lib.GetMaxMessageSize())
		paramBuf := make([]byte, lib.GetMaxMessageSize())
		pctxBuf := make([]byte, lib.GetMaxMessageSize())

	}

}

// this type better implement StoreServer
type myServer struct {
}

func (m *myServer) BestOfAllTime(pctx *parigot.PCtx) (proto.Message, error) {
	out := &pb.BestOfAllTimeResponse{}
	pctx.GetLog()
}

func (m *myServer) Revenue(pctx *parigot.PCtx, in proto.Message) (proto.Message, error) {

}

func (m *myServer) SoldItem(pctx *parigot.PCtx, in proto.Message) error {

}

func register() (string, error) {
	impl := &myServer{}

	_, BestOfAllTimeerr := lib.BindMethodOut(&kernel.BindMethodRequest{
		ProtoPackage: "demo.vvv",
		Service:      "Store",
		Method:       "BestOfAllTime",
	}, impl.BestOfAllTime)
	if BestOfAllTimeerr != nil {
		return "BestOfAllTime", BestOfAllTimeerr
	}

	_, Revenueerr := lib.BindMethodBoth(&kernel.BindMethodRequest{
		ProtoPackage: "demo.vvv",
		Service:      "Store",
		Method:       "Revenue",
	}, impl.Revenue)
	if Revenueerr != nil {
		return "Revenue", Revenueerr
	}

	_, SoldItemerr := lib.BindMethodIn(&kernel.BindMethodRequest{
		ProtoPackage: "demo.vvv",
		Service:      "Store",
		Method:       "SoldItem",
	}, impl.SoldItem)
	if SoldItemerr != nil {
		return "SoldItem", SoldItemerr
	}
	return "", nil
}

func BlockUntilCall(pctx, param, method []byte) {

}

//go:noinline
//go:linkname locate parigot.locate_
func blockUntilCall(int32)

package main

import (
	"fmt"
	"time"
	_ "unsafe"

	"example/vvv/g/business/v1"
	"example/vvv/g/msg/business/v1"

	"github.com/iansmith/parigot/api_impl/syscall"
	"github.com/iansmith/parigot/g/log/v1"
	pblog "github.com/iansmith/parigot/g/msg/log/v1"
	pbsys "github.com/iansmith/parigot/g/msg/syscall/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var logger log.Log
var callImpl = syscall.NewCallImpl()

//go:noinline
func main() {
	//flag.Parse() <--- can't do this until we get startup args figured out

	if _, err := callImpl.Require1("demo.vvv", "Store"); err != nil {
		panic("unable to require my service: " + err.Error())
	}
	if _, err := callImpl.Require1("log", "Log"); err != nil {
		panic("unable to require log service: " + err.Error())
	}
	if _, err := callImpl.Run(&pbsys.RunRequest{Wait: true}); err != nil {
		panic("error starting client process:" + err.Error())
	}
	var err error
	logger, err = log.LocateLog()
	if err != nil {
		Log(pblog.LogLevel_LOG_LEVEL_FATAL, "failed to locate log:%v", err)
	}
	vinnysStore, err := vvv.LocateStore(logger)
	if err != nil {
		Log(pblog.LogLevel_LOG_LEVEL_FATAL, "failed to locate store:%v", err)
	}
	err = vinnysStore.SoldItem(&pb.SoldItemRequest{
		Amount: 14.99,
		When:   timestamppb.New(time.Now()),
	})
	Log(pblog.LogLevel_LOG_LEVEL_INFO, fmt.Sprintf("SoldItem returned ok?:  %v",
		err == nil))
	req := pb.BestOfAllTimeRequest{
		Ctype: pb.ContentType_CONTENT_TYPE_MUSIC,
	}
	//best := &pb.BestOfAllTimeResponse{}
	best, err := vinnysStore.BestOfAllTime(&req)
	if err != nil {
		Log(pblog.LogLevel_LOG_LEVEL_FATAL, "BestOfAllTime failure: %s", err.Error())
	}
	Log(pblog.LogLevel_LOG_LEVEL_INFO, fmt.Sprintf("vinny's BOAT for content %s is: %s, %s, %d", req.Ctype.String(),
		best.Item.Creator, best.Item.Title, best.Item.Year))

	inStock, err := vinnysStore.MediaTypesInStock()
	if err != nil {
		Log(pblog.LogLevel_LOG_LEVEL_ERROR, fmt.Sprintf("MediaTypesInStock() failed  %s", err.Error()))
	} else {
		Log(pblog.LogLevel_LOG_LEVEL_INFO, fmt.Sprintf("MediaTypesInStock: %d different types", len(inStock.InStock)))
		for _, m := range inStock.GetInStock() {
			Log(pblog.LogLevel_LOG_LEVEL_INFO, fmt.Sprintf("%d -> %s", m.Number(), m.String()))
		}
	}
	callImpl.Exit(&pbsys.ExitRequest{Code: 17})
}

func Log(level pblog.LogLevel, spec string, arg ...interface{}) {
	req := &pblog.LogRequest{
		Stamp:   timestamppb.New(time.Now()), // xxx fix me, should be using the kernel
		Level:   level,
		Message: fmt.Sprintf(spec, arg...),
	}
	if logger == nil {
		print("Internal error in storeclient: logger is nil! "+fmt.Sprintf(spec, arg...), "\n")
		return
	}
	if err := logger.Log(req); err != nil {
		print("StoreClient: error in log call:", err.Error(), "\n")
	}
}

func storeclientPrint(spec string, arg ...interface{}) {
	Log(pblog.LogLevel_LOG_LEVEL_DEBUG, "STORECLIENT:"+spec, arg...)
}

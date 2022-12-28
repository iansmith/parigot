package main

import (
	"fmt"
	"time"
	_ "unsafe"

	"demo/vvv/proto/g/vvv"
	"demo/vvv/proto/g/vvv/pb"

	"github.com/iansmith/parigot/api/proto/g/log"
	pblog "github.com/iansmith/parigot/api/proto/g/pb/log"
	"github.com/iansmith/parigot/api/syscall"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var logger log.Log

//go:noinline
func main() {
	//flag.Parse() <--- can't do this until we get startup args figured out

	if _, err := syscall.Require1("demo.vvv", "Store"); err != nil {
		panic("unable to require my service: " + err.Error())
	}
	if _, err := syscall.Require1("log", "Log"); err != nil {
		panic("unable to require log service: " + err.Error())
	}
	if _, err := syscall.Run(true); err != nil {
		panic("error starting client process:" + err.Error())
	}

	vinnysStore, err := vvv.LocateStore()
	if err != nil {
		syscall.Exit(1)
	}
	logger, err = log.LocateLog()
	if err != nil {
		syscall.Exit(1)
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
		storeclientPrint("BestOfAllTime failed %s", err.Error())
		syscall.Exit(1)
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
}

func Log(level pblog.LogLevel, message string) {
	req := &pblog.LogRequest{
		Stamp:   timestamppb.New(time.Now()), // xxx fix me, should be using the kernel
		Level:   level,
		Message: message,
	}
	if err := logger.Log(req); err != nil {
		print("CLIENTSIDESERVICE: error in log call:", err.Error(), "\n")
	}
}

func storeclientPrint(spec string, arg ...interface{}) {
	print("STORECLIENT:", fmt.Sprintf(spec, arg...), "\n")
}

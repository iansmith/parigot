package main

import (
	"time"
	_ "unsafe"

	"demo/vvv/proto/g/vvv"
	"demo/vvv/proto/g/vvv/pb"

	"github.com/iansmith/parigot/g/pb/kernel"
	"github.com/iansmith/parigot/lib"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:noinline
func main() {
	//flag.Parse()
	// logger, err := log.LocateLog()
	// if err != nil {
	// 	print("failed to get log\n")
	// 	//abandon ship, can't get logger to even say what happened
	// 	lib.Exit(&kernel.ExitRequest{Code: 1})
	// }
	// print("xxx trying to set prefix storeclient\n")
	// err = logger.SetPrefix(&pblog.SetPrefixRequest{Prefix: "storeclient"})
	// if err != nil {
	// 	print("xxx set prefix returned error ", err.Error(), "\n")
	// }
	// print("xxx finished set prefix in storeclient", err == nil, "\n")
	// logger.Log(&pblog.LogRequest{Level: pblog.LogLevel_LOGLEVEL_INFO, Message: "Test 1 2 3"})

	vinnysStore, err := vvv.LocateStore()
	if err != nil {
		//logger.Log(&pblog.LogRequest{Level: pblog.LogLevel_LOGLEVEL_FATAL, Message: "could not find the store:" + err.Error()})
		lib.Exit(&kernel.ExitRequest{Code: 1})
	}
	vinnysStore.EnablePctx()

	//t := kernel.Now()
	//logger.LogDebug(fmt.Sprintf("time is now %d ", t), "")
	err = vinnysStore.SoldItem(&pb.SoldItemRequest{
		Amount: 14.99,
		When:   timestamppb.New(time.Now()),
	})
	print("STORECLIENT: got the result for sold item: ", err, "\n")
	req := pb.BestOfAllTimeRequest{
		Ctype: pb.ContentType_CONTENT_TYPE_MUSIC,
	}
	best := &pb.BestOfAllTimeResponse{}
	err = vinnysStore.BestOfAllTime(&req, best)
	if err != nil {
		print("STORE CLIENT, BEST OF ALL TIME ", err.Error(), "\n")
	}
	print("STORECLIENT ", best.Item.Creator, ",", best.Item.Title, "\n")
	// if err != nil {
	//	logger.LogFatal("could not reach the BestOfAllTime call:"+err.Error(), "")
	//}
	//logger.LogDebug("best of all time:"+best.GetMedia().GetTitle(), "")
}

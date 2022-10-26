package main

import (
	"demo/vvv/proto/g/vvv/pb"
	"flag"
	"time"

	"demo/vvv/proto/g/vvv"

	"github.com/iansmith/parigot/g/log"
	"github.com/iansmith/parigot/g/pb/kernel"
	log2 "github.com/iansmith/parigot/g/pb/log"
	"github.com/iansmith/parigot/lib"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:export PkgPathHack
func PkgPathHack(retval int32, x1 int32, x2 int32) {
	return
}

//go:export MethodByNameHack
func MethodByNameHack(retval int32, p0 int32, p1 int32, p2 int32, p3 int32) {
	return
}

//go:export ZeroHack
func ZeroHack(retval int32, p0 int32, p1 int32, p2 int32) {
	return
}

//export main.main
func main() {
	print("hello, fart")
}
func main2() {
	print("1")
	flag.Parse()
	print("2")
	logger, err := log.LocateLog()
	print("3")
	if err != nil {
		//abandon ship, can't get logger to even say what happened
		lib.Exit(&kernel.ExitRequest{Code: 1})
	}
	logger.Log(&log2.LogRequest{Level: 3, Message: "starting up..."})
	vinnysStore, err := vvv.LocateStore()
	if err != nil {
		logger.Log(&log2.LogRequest{Level: 5, Message: "could not find the store:" + err.Error()})
	}
	//t := kernel.Now()
	//logger.LogDebug(fmt.Sprintf("time is now %d ", t), "")
	vinnysStore.SoldItem(&pb.SoldItemRequest{
		Amount: 14.99,
		When:   timestamppb.New(time.Now()),
	})
	//best, err := vinnysStore.BestOfAllTime()
	//if err != nil {
	//	logger.LogFatal("could not reach the BestOfAllTime call:"+err.Error(), "")
	//}
	//logger.LogDebug("best of all time:"+best.GetMedia().GetTitle(), "")
}

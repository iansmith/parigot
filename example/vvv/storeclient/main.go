package main

import (
	"flag"
	"fmt"
	"github.com/iansmith/parigot/g/pb/kernel"
	log2 "github.com/iansmith/parigot/g/pb/log"
	"github.com/iansmith/parigot/lib"
	"time"

	"demo/vvv/proto/g/vvv"

	"github.com/iansmith/parigot/g/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:export PkgPathHack
func PkgPathHack(x2 int32, x3 int32) string {
	return fmt.Sprintf("unknown%d.%d.%d", x2, x3)
}

//export main.main
func main() {
	flag.Parse()
	logger, err := log.LocateLog()
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
	vinnysStore.SoldItem(vvv.SoldItemRequest{
		Amount: 14.99,
		When:   timestamppb.New(time.Now()),
	})
	//best, err := vinnysStore.BestOfAllTime()
	//if err != nil {
	//	logger.LogFatal("could not reach the BestOfAllTime call:"+err.Error(), "")
	//}
	//logger.LogDebug("best of all time:"+best.GetMedia().GetTitle(), "")
}

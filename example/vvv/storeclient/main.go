package main

import (
	"demo/vvv/proto/g/vvv"
	"flag"
	"fmt"
	"github.com/iansmith/parigot/g/parigot/abi"
	"github.com/iansmith/parigot/g/parigot/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

//go:export foo
func foo() {
	abi.Exit(83)
}

//export main.main
func main() {
	flag.Parse()
	logger, err := log.LocateLog()
	if err != nil {
		//abandon ship, can't get logger to even say what happened
		abi.Exit(1)
	}
	logger.LogDebug("starting up", "")
	vinnysStore, err := vvv.LocateStore()
	if err != nil {
		logger.LogFatal("could not find the store:"+err.Error(), "")
	}
	t := abi.Now()
	logger.LogDebug(fmt.Sprintf("time is now %d ", t), "")
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

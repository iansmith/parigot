package main

import (
	"demo/vvv/proto/g/vvv/pb"
	"time"

	"demo/vvv/proto/g/vvv"

	"github.com/iansmith/parigot/g/log"
	"github.com/iansmith/parigot/g/pb/kernel"
	log2 "github.com/iansmith/parigot/g/pb/log"
	"github.com/iansmith/parigot/lib"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	print("xxx1\n")
	//flag.Parse()
	print("xxx2\n")
	logger, err := log.LocateLog()
	print("xxx3\n")
	if err != nil {
		print("failed to get log\n")
		//abandon ship, can't get logger to even say what happened
		lib.Exit(&kernel.ExitRequest{Code: 1})
	}
	print("got log", logger, "\n")
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

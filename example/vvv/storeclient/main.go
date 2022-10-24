package main

import (
	"flag"
	"fmt"
	"time"

	"demo/vvv/proto/g/vvv"

	"github.com/iansmith/parigot/g/parigot/abi"
	"github.com/iansmith/parigot/g/parigot/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:export foo
func foo() {
	abi.Exit(83)
}

//go:export recv
func recv(fn func(blob []byte)) {
	fn([]byte{})
}

//go:export recv64
func recv64(p1 int32, p2 int64) {
	fmt.Sprintf("%d,%d", p1, p2)
}

//go:export PkgPathHack
func PkgPathHack(x2 int32, x3 int32) string {
	return fmt.Sprintf("unknown%d.%d.%d", x2, x3)
}

func Foobie(blob []byte) {

}

//export main.main
func main() {
	flag.Parse()
	recv(Foobie)
	recv64(1, 1023)
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

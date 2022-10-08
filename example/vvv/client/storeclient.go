package main

import (
	//"parigot/api/gen/atlanta1/parigot/log"
	"demo/vvv/proto/gen/demo/vvv"
	"flag"
	"github.com/iansmith/parigot/lib/base/go/log"
)

func main() {
	flag.Parse()
	vinnysStore := vvv.ConnectStore(nil /*params that would be used in prod*/)
	best := vinnysStore.BestOfAllTime()
	log.Dev.Debug("vinny claims the best of all time is ", best.Title)
}

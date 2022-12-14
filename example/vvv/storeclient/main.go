package main

import (
	"fmt"
	"time"
	_ "unsafe"

	"demo/vvv/proto/g/vvv"
	"demo/vvv/proto/g/vvv/pb"

	"github.com/iansmith/parigot/lib"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:noinline
func main() {
	//flag.Parse() <--- can't do this until we get startup args figured out

	print("xxx storeclient 1\n")
	if _, err := lib.Require1("demo.vvv", "Store"); err != nil {
		panic("unable to require my service: " + err.Error())
	}
	print("xxx storeclient 2\n")
	// if _, err := lib.Require1("log", "Log"); err != nil {
	// 	panic("unable to require log service: " + err.Error())
	// }
	print("xxx storeclient 3\n")
	if _, err := lib.Run(true); err != nil {
		panic("error starting client process:" + err.Error())
	}
	print("xxx storeclient 4\n")

	vinnysStore, err := vvv.LocateStore()
	if err != nil {
		lib.Exit(1)
	}
	print("xxx storeclient 5\n")
	// logger, err := log.LocateLog()
	// if err != nil {
	// 	lib.Exit(1)
	// }
	print("xxx storeclient 6\n")
	// logger.Log(&pblog.LogRequest{
	// 	Level:   lib.InfoLevel,
	// 	Message: fmt.Sprintf("About to call sold item implementation..."),
	// })
	print("xxx storeclient 8\n")
	err = vinnysStore.SoldItem(&pb.SoldItemRequest{
		Amount: 14.99,
		When:   timestamppb.New(time.Now()),
	})
	print("xxx storeclient 9\n")
	vinnysStore.Log(lib.InfoLevel, fmt.Sprintf("SoldItem returned ok?:  %v", err == nil))
	print("xxx storeclient 10\n")
	//vinnysStore.DumpLog()
	req := pb.BestOfAllTimeRequest{
		Ctype: pb.ContentType_CONTENT_TYPE_MUSIC,
	}
	//best := &pb.BestOfAllTimeResponse{}
	best, err := vinnysStore.BestOfAllTime(&req)
	if err != nil {
		storeclientPrint("BestOfAllTime failed %s", err.Error())
		lib.Exit(1)
	}
	storeclientPrint("vinny's BOAT for content %s is: %s, %s, %d", req.Ctype.String(),
		best.Item.Creator, best.Item.Title, best.Item.Year)

	inStock, err := vinnysStore.MediaTypesInStock()
	if err != nil {
		storeclientPrint("MediaTypesInStock() failed  %s", err.Error())
	} else {
		storeclientPrint("MediaTypesInStock: %d", len(inStock.InStock))
		print("\t")
		for i, m := range inStock.GetInStock() {
			//print(m.String()) xxxfixme WHYWHYWHY?
			print(m.Number(), " -> ", m.String())
			if i != len(inStock.GetInStock())-1 {
				print(",")
			}
		}
		print("\n")
	}
}

func storeclientPrint(spec string, arg ...interface{}) {
	print("STORECLIENT:", fmt.Sprintf(spec, arg...), "\n")
}

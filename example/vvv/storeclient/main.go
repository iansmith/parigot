package main

import (
	"fmt"
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
	//flag.Parse() <--- can't do this until we get startup args figured out

	vinnysStore, err := vvv.LocateStore()
	if err != nil {
		lib.Exit(&kernel.ExitRequest{Code: 1})
	}
	vinnysStore.EnablePctx()

	err = vinnysStore.SoldItem(&pb.SoldItemRequest{
		Amount: 14.99,
		When:   timestamppb.New(time.Now()),
	})
	storeclientPrint(" SoldItem returned ok?:  %v", err == nil)
	req := pb.BestOfAllTimeRequest{
		Ctype: pb.ContentType_CONTENT_TYPE_MUSIC,
	}
	//best := &pb.BestOfAllTimeResponse{}
	best, err := vinnysStore.BestOfAllTime(&req)
	if err != nil {
		storeclientPrint("BestOfAllTime failed %s", err.Error())
		lib.Exit(&kernel.ExitRequest{Code: 1})
	}
	storeclientPrint("vinny's BOAT for content %s is: %s, %s, %d", req.Ctype.String(),
		best.Item.Creator, best.Item.Title, best.Item.Year)

	inStock, err := vinnysStore.MediaTypesInStock()
	if err != nil {
		storeclientPrint("MediaTypesInStock() failed  %s", err.Error())
	} else {
		storeclientPrint("MediaTypesInStock: %d", len(inStock.InStock))
		print("\t xxx if I print the strings here I get a 'makeslice: len out of range' crash\n")
		print("\t")
		for i, m := range inStock.GetInStock() {
			//print(m.String(), " ") // ARRRRGH
			print(m, " ")
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

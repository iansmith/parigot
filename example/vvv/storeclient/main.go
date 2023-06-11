package main

import (
	"fmt"
	"log"
	_ "unsafe"

	storemsg "example/vvv/g/msg/store/v1"
	"example/vvv/g/store/v1"

	"github.com/iansmith/parigot/apiwasm/syscall"
	lib "github.com/iansmith/parigot/lib/go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:noinline
func main() {
	lib.FlagParseCreateEnv()

	_, err := lib.Require1("demo.vvv", "StoreService")
	if err.IsError() {
		panic("unable to require my service: " + err.Error())
	}
	if _, err := callImpl.Require1("log", "LogService"); err != nil {
		panic("unable to require log service: " + err.Error())
	}
	if _, err := callImpl.Run(&syscall.RunRequest{Wait: true}); err != nil {
		panic("error starting client process:" + err.Error())
	}
	var err error
	vinnysStore, err := store.LocateStoreService(logger)
	if err != nil {
		panic(fmt.Sprintf("failed to locate store:%v", err))
	}
	err = vinnysStore.SoldItem(&storemsg.SoldItemRequest{
		Item: &storemsg.Item{
			Id: 7261991,
		},
		When: timestamppb.Now(), //xxx use kernel Now()
		Amount: &storemsg.Amount{
			Units:      19,
			Hundredths: 99,
		},
	})
	req := storemsg.BestOfAllTimeRequest{
		Content: storemsg.ContentType_CONTENT_TYPE_MUSIC,
	}
	//best := &pb.BestOfAllTimeResponse{}

	_, err = vinnysStore.BestOfAllTime(&req)
	if err != nil {
		panic(fmt.Sprintf("BestOfAllTime failure: %s", err.Error()))
	}
	inStock, err := vinnysStore.MediaTypesInStock()
	if err != nil {
		panic(fmt.Sprintf("MediaTypesInStock() failed  %s", err.Error()))
	} else {
		for _, m := range inStock.GetInStock() {
			log.Printf("%d -> %s", m.Number(), m.String())
		}
	}
	callImpl.Exit(&syscall.ExitRequest{Code: 17})
}

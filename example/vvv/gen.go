package main

import (
	"fmt"

	"demo/vvv/proto/g/vvv/pb"

	"github.com/iansmith/parigot/g/pb/kernel"
	"github.com/iansmith/parigot/lib"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var storeserverVerbose = false

type StoreServer interface {
	BestOfAllTime(pctx lib.Pctx, in proto.Message) (proto.Message, error)
	Revenue(pctx lib.Pctx, in proto.Message) (proto.Message, error)
	SoldItem(pctx lib.Pctx, in proto.Message) error
	MediaTypesInStock(pctx lib.Pctx) (proto.Message, error)
}

var BestOfAllTimeMethod lib.Id
var RevenueMethod lib.Id
var SoldItemMethod lib.Id
var MediaTypesInStockMethod lib.Id

func run(impl StoreServer) {
	// register all methods
	method, err := Storebind()
	if err != nil {
		panic("failed to register method successfully: " + method + ":" + err.Error())
	}
	// allocte buffers for values coming back to use
	paramBuf := make([]byte, lib.GetMaxMessageSize())
	pctxBuf := make([]byte, lib.GetMaxMessageSize())

	// loop on handling calls
	for {
		//
		// wait for notification
		//
		resp, err := StoreBlockUntilCall(pctxBuf, paramBuf)
		if err != nil {
			// error is likely local to this process
			storeserverPrint("RUN:primary for loop ", "Unable to dispatch method call: %v", err)
			continue
		}
		storeserverPrint("RUN: primary for loop ", "block completed, got two values:%d,%d",
			resp.PctxLen, resp.ParamLen)
		//
		// incoming values, pctx and params
		//
		var pctxSlice []byte
		if resp.PctxLen == 0 {
			pctxSlice = []byte{}
		} else {
			pctxSlice = pctxBuf[:resp.PctxLen]
		}
		paramSlice := paramBuf[:resp.ParamLen]
		mid := lib.UnmarshalMethodId(resp.GetMethod())
		cid := lib.UnmarshalCallId(resp.GetCall())

		//
		// create the generic params, pctx and param
		//
		var pctx lib.Pctx
		err = nil
		if resp.PctxLen != 0 {
			pctx, err = lib.NewPctxFromBytes(pctxSlice)
		}
		if err != nil {
			storeserverPrint("RUN: primary for loop ", "Unable to create Pctx for call: %v", err)
			continue
		}
		// a is an any that represents the params
		a := &anypb.Any{}
		err = proto.Unmarshal(paramSlice, a)
		if err != nil {
			storeserverPrint("RUN: primary for loop ", "Unable to create parameters for call: %v", err)
			continue
		}

		//
		// pick the method
		//
		var marshalError, execError error
		var out proto.Message
		switch {
		case mid.Equal(BestOfAllTimeMethod):
			req := &pb.BestOfAllTimeRequest{}
			marshalError = a.UnmarshalTo(req)
			storeserverPrint("RUN:method switch ", "BOAT: unmarshalled request with a ok? [%s]", a.GetTypeUrl())
			if marshalError != nil {
				break
			}
			storeserverPrint("RUN:method switch ", "BOAT: making call")
			out, execError = impl.BestOfAllTime(pctx, req)
			if execError != nil {
				break
			}
			storeserverPrint("RUN: method switch ", "BOAT: got out, ok? %v", out != nil)
		case mid.Equal(RevenueMethod):
			req := &pb.RevenueRequest{}
			marshalError = a.UnmarshalTo(req)
			if marshalError != nil {
				break
			}
			out, execError = impl.Revenue(pctx, req)
			if execError != nil {
				break
			}
		case mid.Equal(SoldItemMethod):
			req := &pb.SoldItemRequest{}
			marshalError = a.UnmarshalTo(req)
			if marshalError != nil {
				break
			}
			execError = impl.SoldItem(pctx, req)
			if execError != nil {
				break
			}
			out = nil // just to be sure
		case mid.Equal(MediaTypesInStockMethod):
			out, execError = impl.MediaTypesInStock(pctx)
			if execError != nil {
				break
			}
		}

		//
		// could be error, could be everything is cool, send to lib to figure it out
		//
		lib.ReturnValueEncode(cid, mid, marshalError, execError, out, pctx)
		// about to loop again
	}
	// wont reach here
}

func Storebind() (string, error) {
	impl := &myServer{}

	resp, BestOfAllTimeerr := lib.BindMethodBoth(&kernel.BindMethodRequest{
		ProtoPackage: "demo.vvv",
		Service:      "Store",
		Method:       "BestOfAllTime",
	}, impl.BestOfAllTime)
	if BestOfAllTimeerr != nil {
		return "BestOfAllTime", BestOfAllTimeerr
	}
	BestOfAllTimeMethod = lib.UnmarshalMethodId(resp.GetMethodId())

	resp, Revenueerr := lib.BindMethodBoth(&kernel.BindMethodRequest{
		ProtoPackage: "demo.vvv",
		Service:      "Store",
		Method:       "Revenue",
	}, impl.Revenue)
	if Revenueerr != nil {
		return "Revenue", Revenueerr
	}
	RevenueMethod = lib.UnmarshalMethodId(resp.GetMethodId())

	resp, SoldItemerr := lib.BindMethodIn(&kernel.BindMethodRequest{
		ProtoPackage: "demo.vvv",
		Service:      "Store",
		Method:       "SoldItem",
	}, impl.SoldItem)
	if SoldItemerr != nil {
		return "SoldItem", SoldItemerr
	}
	SoldItemMethod = lib.UnmarshalMethodId(resp.GetMethodId())

	resp, MediaTypesInStockerr := lib.BindMethodOut(&kernel.BindMethodRequest{
		ProtoPackage: "demo.vvv",
		Service:      "Store",
		Method:       "MediaTypesInStock",
	}, impl.MediaTypesInStock)
	if MediaTypesInStockerr != nil {
		return "MediaTypesInStock", MediaTypesInStockerr
	}
	MediaTypesInStockMethod = lib.UnmarshalMethodId(resp.GetMethodId())
	return "", nil
}

func StoreBlockUntilCall(pctx, param []byte) (*kernel.BlockUntilCallResponse, error) {

	req := &kernel.BlockUntilCallRequest{
		PctxBuffer:  pctx,
		ParamBuffer: param,
	}
	resp, err := lib.BlockUntilCall(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func storeserverPrint(method string, spec string, arg ...interface{}) {
	if storeserverVerbose {
		part1 := fmt.Sprintf("STORESERVER:%s", method)
		part2 := fmt.Sprintf(spec, arg...)
		print(part1, part2, "\n")
	}
}

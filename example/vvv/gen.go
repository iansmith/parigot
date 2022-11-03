package main

import (
	"fmt"
	"log"

	"demo/vvv/proto/g/vvv/pb"

	"github.com/iansmith/parigot/g/pb/kernel"
	"github.com/iansmith/parigot/lib"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type StoreServer interface {
	BestOfAllTime(pctx lib.Pctx, in proto.Message) (proto.Message, error)
	Revenue(pctx lib.Pctx, in proto.Message) (proto.Message, error)
	SoldItem(pctx lib.Pctx, in proto.Message) error
}

var BestOfAllTimeMethod lib.Id
var RevenueMethod lib.Id
var SoldItemMethod lib.Id

func run(impl StoreServer) {
	// register all methods
	method, err := Storebind()
	if err != nil {
		panic("failed to register method successfully: " + method + ":" + err.Error())
	}
	// allocte buffers for values coming back to use
	paramBuf := make([]byte, lib.GetMaxMessageSize())
	pctxBuf := make([]byte, lib.GetMaxMessageSize())

	print("SERVER -- about hit for loop\n")
	// loop on handling calls
	for {
		//
		// wait for notification
		//
		print("STORESERVER about to block\n")
		resp, err := StoreBlockUntilCall(pctxBuf, paramBuf)
		if err != nil {
			// error is likely local to this process
			log.Printf("Unable to dispatch method call: %v", err)
			continue
		}
		print(fmt.Sprintf("STORESERVER block completed, got two values:%d,%d\n",
			resp.PctxLen, resp.ParamLen))
		//
		// incoming values, pctx and params
		//
		pctxSlice := pctxBuf[:resp.PctxLen]
		paramSlice := paramBuf[:resp.ParamLen]
		mid := lib.UnmarshalMethodId(resp.GetMethod())
		cid := lib.UnmarshalCallId(resp.GetCall())

		print(fmt.Sprintf("STORE SERVER unmarshalled arguments %s,%s\n", mid.Short(), cid.Short()))
		//
		// create the generic params, pctx and param
		//
		pctx, err := lib.NewPctxFromBytes(pctxSlice)
		if err != nil {
			log.Printf("Unable to create Pctx for call: %v", err)
			continue
		}
		print(fmt.Sprintf("STORESERVER got pctx\n"))
		// a is an any that represents the params
		a := &anypb.Any{}
		err = proto.Unmarshal(paramSlice, a)
		if err != nil {
			log.Printf("Unable to create parameters for call: %v", err)
			continue
		}
		print(fmt.Sprintf("STORESERVER got param\n"))

		//
		// pick the method
		//
		var marshalError, execError error
		var out proto.Message
		switch {
		case mid.Equal(BestOfAllTimeMethod):
			req := &pb.BestOfAllTimeRequest{}
			marshalError = a.UnmarshalTo(req)
			if marshalError != nil {
				break
			}
			out, execError = impl.BestOfAllTime(pctx, req)
			if execError != nil {
				break
			}
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
			print("STORESERVER got a call on sold item\n")
			req := &pb.SoldItemRequest{}
			marshalError = a.UnmarshalTo(req)
			if marshalError != nil {
				break
			}
			execError = impl.SoldItem(pctx, req)
			if execError != nil {
				break
			}
			print("STORESERVER executed call to sold item", execError == nil, "\n")
			out = nil // just to be sure
		}

		print("STORESERVER executed call to sold item 1111\n")
		// now check for errors
		if marshalError != nil {
			log.Printf("unable to unmarshal parameters:%v", marshalError)
			lib.ReturnValue(&kernel.ReturnValueRequest{
				Call:         lib.MarshalCallId(cid),
				ErrorMessage: marshalError.Error(),
			})
			continue
		}
		print("STORESERVER executed call to sold item 2222\n")
		// xxx fixme: this should be doing something special here to indicate the USER code barfed
		if execError != nil {
			log.Printf("STORESERVER unable to exec user function:%v", execError)
			lib.ReturnValue(&kernel.ReturnValueRequest{
				Call:         lib.MarshalCallId(cid),
				ErrorMessage: execError.Error(),
			})
			continue
		}
		print("STORESERVER executed result aappears to be ok\n")

		//
		// everything is cool, send results
		//
		var resultBuf []byte
		var marshalErr error
		if out != nil {
			resultBuf, err = proto.Marshal(out)
		} else {
			resultBuf = nil
		}
		// result problem?
		if marshalErr != nil {
			print(fmt.Sprintf("STORESERVER unable to marshal results:%v\n", marshalErr))
			lib.ReturnValue(&kernel.ReturnValueRequest{
				Call:         lib.MarshalCallId(cid),
				ErrorMessage: execError.Error(),
			})
		}

		print("STORESERVER sending response???\n")
		// success!
		lib.ReturnValue(&kernel.ReturnValueRequest{
			Call:         lib.MarshalCallId(cid),
			ResultBuffer: resultBuf,
		})

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
	print("CLIENT best of all time: ", BestOfAllTimeMethod.Short()+"\n")

	resp, Revenueerr := lib.BindMethodBoth(&kernel.BindMethodRequest{
		ProtoPackage: "demo.vvv",
		Service:      "Store",
		Method:       "Revenue",
	}, impl.Revenue)
	if Revenueerr != nil {
		return "Revenue", Revenueerr
	}
	RevenueMethod = lib.UnmarshalMethodId(resp.GetMethodId())
	print("CLIENT revenue: ", RevenueMethod.Short()+"\n")

	resp, SoldItemerr := lib.BindMethodIn(&kernel.BindMethodRequest{
		ProtoPackage: "demo.vvv",
		Service:      "Store",
		Method:       "SoldItem",
	}, impl.SoldItem)
	if SoldItemerr != nil {
		return "SoldItem", SoldItemerr
	}
	SoldItemMethod = lib.UnmarshalMethodId(resp.GetMethodId())
	print("CLIENT sold: ", SoldItemMethod.Short()+"\n")
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

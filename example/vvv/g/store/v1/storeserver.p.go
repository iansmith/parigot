// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: store/v1/store.proto

package store




import (
	"fmt"
storemsg "example/vvv/g/msg/store/v1" 

	os "os"
    // this set of imports is _unrelated_ to the particulars of what the .proto imported... those are above
	"github.com/iansmith/parigot/g/protosupport/v1"
	"github.com/iansmith/parigot/g/syscall/v1"
	"github.com/iansmith/parigot/apiwasm/syscall"
	"github.com/iansmith/parigot/apiwasm/background"
	lib "github.com/iansmith/parigot/lib/go"
	
	"google.golang.org/protobuf/proto"
	
)

//
// StoreServiceServer (from store/v1/store.proto)
//

type StoreServiceServer interface {

	MediaTypesInStock(pctx *protosupport.Pctx)(proto.Message, error) 
	BestOfAllTime(pctx *protosupport.Pctx, in proto.Message)(proto.Message, error) 
	Revenue(pctx *protosupport.Pctx, in proto.Message)(proto.Message, error) 
	SoldItem(pctx *protosupport.Pctx, in proto.Message)error  
	Ready() bool
} 

//
// StoreService method ids
//
var mediaTypesInStockMethod lib.Id
var bestOfAllTimeMethod lib.Id
var revenueMethod lib.Id
var soldItemMethod lib.Id 


var storeServiceServerVerbose = true
var storeServiceCall = syscall.NewCallImpl()
var storeServiceBackground background.Background

func RunStoreService(impl StoreServiceServer) {
	// register all methods
	method, err := storeServiceBind(impl)
	if err != nil {
		panic("failed to register method successfully: " + method + ":" + err.Error())
	}
	// loop on handling calls
	for {
		//
		// wait for notification
		//
		resp, err := storeServiceBlockUntilCall(storeServiceBackground!=nil)
		if err != nil {
			// error is likely local to this process
			storeServicePrint("RUN:primary for loop ", "Unable to dispatch method call: %v", err)
			continue
		}
		if resp.TimedOut{
			if storeServiceBackground==nil {
				continue
			}
			storeServiceBackground.Background()
			continue
		}
		//storeServicePrint("RUN: primary for loop ", "block completed, got two values:pctx size %d, param size %d",
		//	proto.Size(resp.GetPctx()), proto.Size(resp.GetParam()))
		//
		// incoming values, pctx and params
		//
		mid := lib.Unmarshal(resp.GetMethod())
		cid := lib.Unmarshal(resp.GetCall())

		//
		// pick the method
		//
		var marshalError, execError error
		var out proto.Message
		switch {
		case mid.Equal(mediaTypesInStockMethod):
            // no input for this method
			out, execError = impl.MediaTypesInStock(resp.GetPctx()) 
			if execError != nil {
				break
			}
		case mid.Equal(bestOfAllTimeMethod):
			req := &storemsg.BestOfAllTimeRequest{}
			marshalError = resp.GetParam().UnmarshalTo(req)
			if marshalError != nil {
				break
			}
			out, execError = impl.BestOfAllTime(resp.GetPctx(),req) 
			if execError != nil {
				break
			}
		case mid.Equal(revenueMethod):
			req := &storemsg.RevenueRequest{}
			marshalError = resp.GetParam().UnmarshalTo(req)
			if marshalError != nil {
				break
			}
			out, execError = impl.Revenue(resp.GetPctx(),req) 
			if execError != nil {
				break
			}
		case mid.Equal(soldItemMethod):
			req := &storemsg.SoldItemRequest{}
			marshalError = resp.GetParam().UnmarshalTo(req)
			if marshalError != nil {
				break
			}

            execError = impl.SoldItem(resp.GetPctx(), req)
			if execError != nil {
				break
			} 
        }
		//
		// could be error, could be everything is cool, send to lib to figure it out
		//
		retresp, e:=lib.ReturnValueEncode(storeServiceCall, cid, mid, marshalError, execError, out, resp.GetPctx())
		if e!=nil {
			panic("server failed to return value: "+e.Error())
		}
		if retresp.ExitAfterUse {
			os.Exit(2)
		}
		// about to loop again
	}
	// wont reach here
}


func storeServiceBind(impl StoreServiceServer) (string, error) {

	resp, mediaTypesInStockerr := storeServiceCall.BindMethodOut(&syscallmsg.BindMethodRequest{
		ProtoPackage: "store.v1",
		Service:      "StoreService",
		Method:       "MediaTypesInStock",
	}, impl.MediaTypesInStock)
	if mediaTypesInStockerr != nil {
		return "MediaTypesInStock", mediaTypesInStockerr
	}
	mediaTypesInStockMethod = lib.Unmarshal(resp.GetMethodId())

	resp, bestOfAllTimeerr := storeServiceCall.BindMethodBoth(&syscallmsg.BindMethodRequest{
		ProtoPackage: "store.v1",
		Service:      "StoreService",
		Method:       "BestOfAllTime",
	}, impl.BestOfAllTime)
	if bestOfAllTimeerr != nil {
		return "BestOfAllTime", bestOfAllTimeerr
	}
	bestOfAllTimeMethod = lib.Unmarshal(resp.GetMethodId())

	resp, revenueerr := storeServiceCall.BindMethodBoth(&syscallmsg.BindMethodRequest{
		ProtoPackage: "store.v1",
		Service:      "StoreService",
		Method:       "Revenue",
	}, impl.Revenue)
	if revenueerr != nil {
		return "Revenue", revenueerr
	}
	revenueMethod = lib.Unmarshal(resp.GetMethodId())

	resp, soldItemerr := storeServiceCall.BindMethodIn(&syscallmsg.BindMethodRequest{
		ProtoPackage: "store.v1",
		Service:      "StoreService",
		Method:       "SoldItem",
	}, impl.SoldItem)
	if soldItemerr != nil {
		return "SoldItem", soldItemerr
	}
	soldItemMethod = lib.Unmarshal(resp.GetMethodId()) 
	if !impl.Ready(){
		panic("unable to start StoreService because it failed Ready() check")
	}
	var ok bool
	storeServiceBackground, ok = impl.(background.Background)
	if !ok {
		print("StoreServiceServer not background\n")
	}
	return "",nil
}


func storeServiceTimedOut()  {
	print("storeServiceTimedOut()\n")
}

func storeServiceBlockUntilCall(canTimeout bool) (*syscallmsg.BlockUntilCallResponse, error) {
	req := &syscallmsg.BlockUntilCallRequest{}
	resp, err := storeServiceCall.BlockUntilCall(req, canTimeout)
	if err != nil {
		return nil, err
	}
	
	return resp, nil
}

func storeServicePrint(method string, spec string, arg ...interface{}) {
	if storeServiceServerVerbose {
		part1 := fmt.Sprintf("storeServiceServer:%s", method)
		part2 := fmt.Sprintf(spec, arg...)
		print(part1, part2, "\n")
	}
} 

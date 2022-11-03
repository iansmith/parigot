package main

import (
	"fmt"
	"log"

	"github.com/iansmith/parigot/g/pb/kernel"
	pblog "github.com/iansmith/parigot/g/pb/log"
	"github.com/iansmith/parigot/lib"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type TerminalLog interface {
	// xxx we don't have a way to do this signature right now without pctx
	Log(in proto.Message) error
	SetPrefix(in proto.Message) error
}

var LogMethod lib.Id
var SetPrefixMethod lib.Id

func run(impl TerminalLog) {
	print("xxxx terminal log 1111\n")
	// register all methods
	method, err := TerminalLogbind()
	if err != nil {
		panic("failed to register method successfully: " + method + ":" + err.Error())
	}
	print("xxxx terminal log 2222\n")
	// allocte buffers for values coming back to use
	paramBuf := []byte{}
	pctxBuf := make([]byte, lib.GetMaxMessageSize())

	print("simple log, heading for the loop\n")
	print("xxxx terminal log 333\n")
	// loop on handling calls
	for {
		//
		// wait for notification
		//
		print("xxxx terminal log 4444\n")
		print("terminallog -- about to block\n")
		resp, err := TerminalLogBlockUntilCall(pctxBuf, paramBuf)
		if err != nil {
			// error is likely local to this process
			log.Printf("Unable to dispatch method call: %v", err)
			continue
		}
		print("TERMINAL %+v, %+v, %+v\n\n")
		print("xxxx terminal log finished the block 5555, len is ", resp.ParamLen, "\n")
		//
		// incoming values, pctx and params
		//
		// xxx no way to do this right now and dodge pctxSlice
		//pctxSlice := pctxBuf[:resp.PctxLen]
		paramSlice := paramBuf[:resp.ParamLen]
		mid := lib.UnmarshalMethodId(resp.GetMethod())
		cid := lib.UnmarshalCallId(resp.GetCall())
		print(fmt.Sprintf("xxxx terminal log 6666, slicelen=%d %s,%s\n", len(paramSlice), mid.Short(), cid.Short()))

		//
		// create the generic params, pctx and param
		//
		// xxx again, we don't have a way right now to avoid the pctx
		// a is an any that represents the params
		//
		print("xxxx terminal log 666aaaa", mid.Short(), "\n")
		err = proto.Unmarshal(paramSlice, &anypb.Any{})
		print("xxxx terminal log 77777", mid.Short(), "\n")
		if err != nil {
			log.Printf("Unable to create parameters for call: %v", err)
			continue
		}
		a := &anypb.Any{}
		//
		// pick the method
		//
		var marshalError, execError error
		var out proto.Message
		switch {
		case mid.Equal(LogMethod):
			print("xxxx terminal log 8888\n")
			req := &pblog.LogRequest{}
			marshalError = a.UnmarshalTo(req)
			if marshalError != nil {
				break
			}
			execError = impl.Log(req)
		case mid.Equal(SetPrefixMethod):
			print("xxxx terminal log 9999\n")
			req := &pblog.SetPrefixRequest{}
			marshalError = a.UnmarshalTo(req)
			if marshalError != nil {
				break
			}
			execError = impl.SetPrefix(req)
		}
		print("xxxx terminal log 000000\n")
		// now check for errors
		if marshalError != nil {
			log.Printf("unable to unmarshal parameters:%v", marshalError)
			lib.ReturnValue(&kernel.ReturnValueRequest{
				Call:         lib.MarshalCallId(cid),
				ErrorMessage: marshalError.Error(),
			})
			continue
		}
		// xxx fixme: this should be doing something special here to indicate the USER code barfed
		if execError != nil {
			log.Printf("unable to exec user function:%v", execError)
			lib.ReturnValue(&kernel.ReturnValueRequest{
				Call:         lib.MarshalCallId(cid),
				ErrorMessage: execError.Error(),
			})
			continue
		}
		print("terminallog finished execution of ", mid.Short(), "\n")
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
			log.Printf("unable to marshal results:%v", marshalErr)
			lib.ReturnValue(&kernel.ReturnValueRequest{
				Call:         lib.MarshalCallId(cid),
				ErrorMessage: execError.Error(),
			})
		}

		// success!
		lib.ReturnValue(&kernel.ReturnValueRequest{
			Call:         lib.MarshalCallId(cid),
			ResultBuffer: resultBuf,
		})

		// about to loop again
	}
	// wont reach here
}

func TerminalLogbind() (string, error) {
	impl := &terminalLog{}

	print("xxxx terminal logbind 1111\n")
	resp, Logerr := lib.BindMethodInNoPctx(&kernel.BindMethodRequest{
		ProtoPackage: "log",
		Service:      "Log",
		Method:       "Log",
	}, impl.Log)
	if Logerr != nil {
		return "Log", Logerr
	}
	LogMethod = lib.UnmarshalMethodId(resp.GetMethodId())
	print("Log.Log ", LogMethod.Short()+"\n")
	print("xxxx terminal logbind 222\n")
	resp, SetPrefixerr := lib.BindMethodInNoPctx(&kernel.BindMethodRequest{
		ProtoPackage: "log",
		Service:      "Log",
		Method:       "SetPrefix",
	}, impl.SetPrefix)
	if SetPrefixerr != nil {
		return "SetPrefix", SetPrefixerr
	}
	SetPrefixMethod = lib.UnmarshalMethodId(resp.GetMethodId())
	print("SetPrefix.SetPrefix: ", LogMethod.Short()+"\n")

	return "", nil
}

func TerminalLogBlockUntilCall(pctx, param []byte) (*kernel.BlockUntilCallResponse, error) {
	print("TerminalBlockUntilCall: about to block\n")
	req := &kernel.BlockUntilCallRequest{
		PctxBuffer:  pctx,
		ParamBuffer: param,
	}
	resp, err := lib.BlockUntilCall(req)
	print("TerminalBlockUntilCall: returned:", resp == nil, " and ", err == nil, "\n")
	if err != nil {
		print("terminalLogBlockUntilCall, aborting because BlockUntilCall failed: %s\n",
			err.Error())
		return nil, err
	}
	cid := lib.UnmarshalCallId(resp.Call)
	mid := lib.UnmarshalMethodId(resp.Method)
	e := lib.UnmarshalKernelErrorId(resp.ErrorId)
	print("!!TerminalLogBlockUntilCall:", fmt.Sprintf("raw id info from resp: %s,%s,%s ", cid.Short(),
		mid.Short(), e.Short()), "\n")
	print("!!TerminalLogBlockUntilCall:", fmt.Sprintf("paramlen=%d, pctxlen=%d", resp.ParamLen,
		resp.PctxLen), "\n")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

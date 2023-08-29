//go:build wasip1 

// Code generated by protoc-gen-parigot. DO NOT EDIT.
// source: frontdoor/v1/frontdoor.proto

package frontdoor


import(
    "context" 

// no method? true

    "github.com/iansmith/parigot/lib/go/future"  
    "github.com/iansmith/parigot/lib/go/client"  
    "github.com/iansmith/parigot/api/shared/id"
    "google.golang.org/protobuf/proto"
    syscallguest "github.com/iansmith/parigot/api/guest/syscall" 
    syscall "github.com/iansmith/parigot/g/syscall/v1" 
    "github.com/iansmith/parigot/lib/go"  
    "google.golang.org/protobuf/types/known/anypb"


)  
//
// Frontdoor from frontdoor/v1/frontdoor.proto
//
//service interface
type Frontdoor interface { 
    Ready(context.Context,id.ServiceId) *future.Base[bool]
}

type Client interface { 
}

// Client difference from Frontdoor: Ready() 
type Client_ struct {
    *client.BaseService
}
// Check that Client_ is a Client.
var _ = Client(&Client_{})  
package log
import(
    "github.com/iansmith/parigot/api/proto/g/pb/log" 

    // this set of imports is _unrelated_ to the particulars of what the .proto imported... those are above
    "github.com/iansmith/parigot/lib"  // id manipulation
    "github.com/iansmith/parigot/api/proto/g/pb/call" // dispatch and locate
    "github.com/iansmith/parigot/api/proto/g/pb/protosupport" // ids
)
//
// Log
//
type Log interface {
    Log(in *log.LogRequest)error 
} 



type LogClient struct {
    *lib.ClientSideService
}

func LocateLog() (*LogClient,error) {
	var resp *call.LocateResponse
	req := &call.LocateRequest{
        PackageName:"log",
        ServiceName: "Log",
	}
    resp, err:=lib.CallConnection().Locate(req)
    if err!=nil {
        return nil, err
    }
    service:=lib.Unmarshal[*protosupport.ServiceId](resp.GetServiceId())
    cs := lib.NewClientSideService(service, "LogClient")
    return &LogClient{
        ClientSideService: cs,
    }, nil
}

func (i *LogClient) Log(in *log.LogRequest)error { 
    _, err:= i.Dispatch("Log",in)
 

    if err!=nil {
        return err
    }
  
    return nil 
}
   


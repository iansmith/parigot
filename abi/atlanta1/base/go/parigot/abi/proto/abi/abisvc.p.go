package %s
net
type ParigotAbi interface {
    OutputString(OutputStringRequest)OutputStringRequestExit(ExitRequest)ExitRequestNow()SetNow(SetNowRequest)SetNowRequestTinygoNotImplemented(NotImplementedRequest)NotImplementedRequestJSNotImplemented(NotImplementedRequest)NotImplementedRequest
}

type ParigotAbiHandler struct {
    OutputString func(OutputStringRequest) OutputStringRequest
    Exit func(ExitRequest) ExitRequest
    Now func() 
    SetNow func(SetNowRequest) SetNowRequest
    TinygoNotImplemented func(NotImplementedRequest) NotImplementedRequest
    JSNotImplemented func(NotImplementedRequest) NotImplementedRequest
}

func LocateParigotAbi() parigot_abi {
    anything:=parigot.Locate("parigot_abi")
    return anything.(ParigotAbi)
}

func RegisterParigotAbi(h parigot_abiHandler) {
    parigot.Register("parigot_abi",h)
}


package main

import (
	"fmt"
	"time"

	"demo/vvv/proto/g/vvv"
	"demo/vvv/proto/g/vvv/pb"

	"github.com/iansmith/parigot/api/proto/g/file"
	"github.com/iansmith/parigot/api/proto/g/log"
	"github.com/iansmith/parigot/api/proto/g/pb/call"
	pbfile "github.com/iansmith/parigot/api/proto/g/pb/file"
	pblog "github.com/iansmith/parigot/api/proto/g/pb/log"
	"github.com/iansmith/parigot/api/proto/g/pb/protosupport"
	"github.com/iansmith/parigot/api/syscall"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var callImpl = syscall.NewCallImpl()

func main() {
	//if things need to be required/exported you need to force them to the ready state BEFORE calling run()
	if _, err := callImpl.Require1("log", "Log"); err != nil {
		panic("unable to require log service: " + err.Error())
	}
	if _, err := callImpl.Require1("file", "File"); err != nil {
		panic("unable to require file service:" + err.Error())
	}
	if _, err := callImpl.Export1("demo.vvv", "Store"); err != nil {
		panic("unable to export demo.vvv: " + err.Error())
	}
	vvv.Run(&myServer{})
}

// this type better implement vvv.StoreServer
type myServer struct {
	logger  log.Log
	fileSvc file.File
}

//
// This file contains the true implementations--the server side--for the methods
// defined in business.proto.
//

func (m *myServer) MediaTypesInStock(pctx *protosupport.Pctx) (proto.Message, error) {
	return &pb.MediaTypesInStockResponse{
		InStock: []pb.Media{pb.Media_MEDIA_BETA, pb.Media_MEDIA_CASSETTE, pb.Media_MEDIA_LASER_DISC, pb.Media_MEDIA_VHS, pb.Media_MEDIA_VINYL, pb.Media_MEDIA_ATARI_CART},
	}, nil
}

func (m *myServer) BestOfAllTime(pctx *protosupport.Pctx, inProto proto.Message) (proto.Message, error) {
	in := inProto.(*pb.BestOfAllTimeRequest)

	out := &pb.BestOfAllTimeResponse{
		Item: &pb.Item{},
	}
	m.log(pctx, pblog.LogLevel_LOG_LEVEL_DEBUG, "reached BestOfAllTime, computing his choices")

	if in.Ctype == pb.ContentType_CONTENT_TYPE_MUSIC {
		out.Item.Creator = "The Smiths"
		out.Item.Title = "The Queen Is Dead"
		out.Item.Year = 1984
		out.Item.Ctype = pb.ContentType_CONTENT_TYPE_MUSIC
		out.Item.Media = pb.Media_MEDIA_CASSETTE
		out.Item.Price = 19.99
		return out, nil
	}
	if in.Ctype == pb.ContentType_CONTENT_TYPE_MOVIE {
		out.Item.Creator = "George Lucas"
		out.Item.Title = "Star Wars Episode IV"
		out.Item.Year = 1979
		out.Item.Ctype = pb.ContentType_CONTENT_TYPE_MOVIE
		out.Item.Media = pb.Media_MEDIA_LASER_DISC
		out.Item.Price = 89.99
		return out, nil
	}
	if in.Ctype == pb.ContentType_CONTENT_TYPE_TV {
		out.Item.Creator = "Columbia House"
		out.Item.Title = "M*A*S*H Collectors Edition"
		out.Item.Year = 1992
		out.Item.Ctype = pb.ContentType_CONTENT_TYPE_TV
		out.Item.Media = pb.Media_MEDIA_VHS
		out.Item.Price = 29.99
		return out, nil
	}
	m.log(pctx, pblog.LogLevel_LOG_LEVEL_INFO, "unexpected content type in request %d", int32(in.Ctype))
	return nil, fmt.Errorf("unexpected content type request int %d", int32(in.Ctype))
}

func (m *myServer) Revenue(pctx *protosupport.Pctx, in proto.Message) (proto.Message, error) {
	out := &pb.RevenueResponse{}
	m.log(pctx, pblog.LogLevel_LOG_LEVEL_WARNING, "Revenue() not yet implemented, ignoring input value, returning dummy values")
	out.Revenue = 817.71
	return out, nil
}

func (m *myServer) SoldItem(pctx *protosupport.Pctx, in proto.Message) error {
	m.log(pctx, pblog.LogLevel_LOG_LEVEL_WARNING, "SoldItem() not yet implemented, ignoring input value")
	return nil
}

// Ready is a check, if this returns false the library should abort and not attempt to run this service.
// Normally, this is used to block using the lib.Run() call.  This call will wait until all the required
// services are ready.
func (m *myServer) Ready() bool {

	if _, err := callImpl.Run(&call.RunRequest{Wait: true}); err != nil {
		print("ready: error in attempt to signal Run: ", err.Error(), "\n")
		return false
	}
	logger, err := log.LocateLog()
	if err != nil {
		print("ERROR trying to create log client: ", err.Error(), "\n")
		return false
	}
	m.logger = logger

	m.log(nil, pblog.LogLevel_LOG_LEVEL_DEBUG, "about to locate file")
	fs, err := file.LocateFile(logger)
	m.log(nil, pblog.LogLevel_LOG_LEVEL_DEBUG, "server located (fs!=nil) %v,(err==nil) %v",
		fs != nil, err == nil)
	if err != nil {
		print("ERROR trying to create fs client: ", err.Error(), "\n")
		return false
	}
	m.fileSvc = fs
	// load the test data
	_, err = m.fileSvc.Load(&pbfile.LoadRequest{Path: "testdata/vvv"})
	if err != nil {
		m.log(nil, pblog.LogLevel_LOG_LEVEL_FATAL, "unable to load the test data: %v: ", err.Error())
	}
	m.log(nil, pblog.LogLevel_LOG_LEVEL_DEBUG, "about to open the file %s", "/app/testdata/vvv/boat.toml")
	resp, err := m.fileSvc.Open(&pbfile.OpenRequest{Path: "/app/testdata/vvv/boat.toml"})
	m.log(nil, pblog.LogLevel_LOG_LEVEL_DEBUG, "open file returned err nil? %v", err == nil)
	if err != nil {
		m.log(nil, pblog.LogLevel_LOG_LEVEL_FATAL, "unable to open the boat.toml file:%v", err)
	}
	m.log(nil, pblog.LogLevel_LOG_LEVEL_DEBUG, "file open ok %s", resp.GetId().String())

	return true

}

func (m *myServer) log(pctx *protosupport.Pctx, level pblog.LogLevel, spec string, rest ...interface{}) {
	n := time.Now()
	if pctx != nil && !pctx.GetNow().AsTime().IsZero() {
		n = pctx.GetNow().AsTime()
	}
	msg := fmt.Sprintf(spec, rest...)
	req := pblog.LogRequest{
		Stamp:   timestamppb.New(n),
		Level:   level,
		Message: msg,
	}
	m.logger.Log(&req)
}

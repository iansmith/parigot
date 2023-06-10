package main

import (
	"fmt"
	"time"

	storemsg "example/vvv/g/msg/store/v1"
	"example/vvv/g/store/v1"

	"github.com/iansmith/parigot/apiimpl/syscall"
	"github.com/iansmith/parigot/g/file/v1"
	"github.com/iansmith/parigot/g/log/v1"
	filemsg "github.com/iansmith/parigot/g/msg/file/v1"
	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	pbsys "github.com/iansmith/parigot/g/msg/syscall/v1"
	lib "github.com/iansmith/parigot/lib/go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var callImpl = syscall.NewCallImpl()

func main() {
	lib.FlagParseCreateEnv()

	//if things need to be required/exported you need to force them to the ready state BEFORE calling run()
	if _, err := callImpl.Require1("log", "LogService"); err != nil {
		panic("unable to require log service: " + err.Error())
	}
	if _, err := callImpl.Require1("file", "FileService"); err != nil {
		panic("unable to require file service:" + err.Error())
	}
	if _, err := callImpl.Export1("demo.vvv", "StoreService"); err != nil {
		panic("unable to export demo.vvv: " + err.Error())
	}
	store.RunStoreService(&myStoreServer{})
}

// this type better implement store.StoreServer
var checkCompat = &myStoreServer{}

type myStoreServer struct {
	logger  log.LogService
	fileSvc file.FileService
}

//
// This file contains the true implementations--the server side--for the methods
// defined in business.proto.
//

func (m *myStoreServer) MediaTypesInStock(pctx *protosupportmsg.Pctx) (proto.Message, error) {
	return &storemsg.MediaTypesInStockResponse{
		InStock: []storemsg.MediaType{
			storemsg.MediaType_MEDIA_TYPE_ATARI_CART,
			storemsg.MediaType_MEDIA_TYPE_BETA,
			storemsg.MediaType_MEDIA_TYPE_CASSETTE,
			storemsg.MediaType_MEDIA_TYPE_LASERDISC,
			storemsg.MediaType_MEDIA_TYPE_VHS,
			storemsg.MediaType_MEDIA_TYPE_VINYL,
			storemsg.MediaType_MEDIA_TYPE_ATARI_CART},
	}, nil
}

func (m *myStoreServer) BestOfAllTime(pctx *protosupportmsg.Pctx, inProto proto.Message) (proto.Message, error) {
	in := inProto.(*storemsg.BestOfAllTimeRequest)

	out := &storemsg.BestOfAllTimeResponse{
		Boat: &storemsg.Boat{},
	}
	m.log(pctx, logmsg.LogLevel_LOG_LEVEL_DEBUG, "reached BestOfAllTime, computing his choices")

	if in.Content == storemsg.ContentType_CONTENT_TYPE_MUSIC {
		out.Boat.Creator = "The Smiths"
		out.Boat.Title = "The Queen Is Dead"
		out.Boat.Year = 1984
		out.Boat.Content = storemsg.ContentType_CONTENT_TYPE_MUSIC
		out.Boat.Media = storemsg.MediaType_MEDIA_TYPE_CASSETTE
		out.Boat.Price = &storemsg.Amount{
			Units:      14,
			Hundredths: 99,
		}
		return out, nil
	}
	// if in.Ctype == storemsg.ContentType_CONTENT_TYPE_MOVIE {
	// 	out.Item.Creator = "George Lucas"
	// 	out.Item.Title = "Star Wars Episode IV"
	// 	out.Item.Year = 1979
	// 	out.Item.Ctype = storemsg.ContentType_CONTENT_TYPE_MOVIE
	// 	out.Item.Media = storemsg.Media_MEDIA_LASER_DISC
	// 	out.Item.Price = 89.99
	// 	return out, nil
	// }
	// if in.Ctype == storemsg.ContentType_CONTENT_TYPE_TV {
	// 	out.Item.Creator = "Columbia House"
	// 	out.Item.Title = "M*A*S*H Collectors Edition"
	// 	out.Item.Year = 1992
	// 	out.Item.Ctype = storemsg.ContentType_CONTENT_TYPE_TV
	// 	out.Item.Media = storemsg.Media_MEDIA_VHS
	// 	out.Item.Price = 29.99
	// 	return out, nil
	// }
	m.log(pctx, logmsg.LogLevel_LOG_LEVEL_INFO, "unexpected content type in request %d", int32(in.Content))
	return nil, fmt.Errorf("unexpected content type request int %d", int32(in.Content))
}

func (m *myStoreServer) Revenue(pctx *protosupportmsg.Pctx, in proto.Message) (proto.Message, error) {
	out := &storemsg.RevenueResponse{}
	m.log(pctx, logmsg.LogLevel_LOG_LEVEL_WARNING, "Revenue() not yet implemented, ignoring input value, returning dummy values")
	out.Revenue = 817.71
	return out, nil
}

func (m *myStoreServer) SoldItem(pctx *protosupportmsg.Pctx, in proto.Message) error {
	m.log(pctx, logmsg.LogLevel_LOG_LEVEL_WARNING, "SoldItem() not yet implemented, ignoring input value")
	return nil
}

// Ready is a check, if this returns false the library should abort and not attempt to run this service.
// Normally, this is used to block using the lib.Run() call.  This call will wait until all the required
// services are ready.
func (m *myStoreServer) Ready() bool {
	if _, err := callImpl.Run(&pbsys.RunRequest{Wait: true}); err != nil {
		print("ready: error in attempt to signal Run: ", err.Error(), "\n")
		return false
	}
	logger, err := log.LocateLogService()
	if err != nil {
		print("ERROR trying to create log client: ", err.Error(), "\n")
		return false
	}
	m.logger = logger

	fClient, err := file.LocateFileService(logger)
	if err != nil {
		print("ERROR trying to create fs client: ", err.Error(), "\n")
		return false
	}
	m.fileSvc = fClient
	// load the test data... this creates an in-memory filesystem that has an entry
	// at /app/testdata/vvv
	_, err = m.fileSvc.LoadTest(&filemsg.LoadTestRequest{
		Path:          "testdata/vvv",
		MountLocation: "/testdata/vvv", // => /app/testdata/vvv/<blah>.toml
	})
	if err != nil {
		m.log(nil, logmsg.LogLevel_LOG_LEVEL_FATAL, "unable to load the test data: %v: ", err.Error())
		return false
	}
	_, err = m.fileSvc.Open(&filemsg.OpenRequest{Path: "/app/testdata/vvv/boat.toml"})
	if err != nil {
		m.log(nil, logmsg.LogLevel_LOG_LEVEL_FATAL, "unable to open the boat.toml file:%v", err)
		return false
	}
	return true

}

func (m *myStoreServer) log(pctx *protosupportmsg.Pctx, level logmsg.LogLevel, spec string, rest ...interface{}) {
	n := time.Now()
	if pctx != nil && !pctx.GetNow().AsTime().IsZero() {
		n = pctx.GetNow().AsTime()
	}
	msg := fmt.Sprintf(spec, rest...)
	req := logmsg.LogRequest{
		Stamp:   timestamppb.New(n),
		Level:   level,
		Message: msg,
	}
	m.logger.Log(&req)
}

package main

import (
	"fmt"

	"demo/vvv/proto/g/vvv"
	"demo/vvv/proto/g/vvv/pb"

	pblog "github.com/iansmith/parigot/api/proto/g/pb/log"
	"github.com/iansmith/parigot/lib"
	"google.golang.org/protobuf/proto"
)

func main() {
	// if things need to be required you need to force them to the ready state BEFORE calling run()
	// if _, err := lib.Require1("log", "Log"); err != nil {
	// 	panic("unable to require log service: " + err.Error())
	// }
	vvv.Run(&myServer{})
}

// this type better implement vvv.StoreServer
type myServer struct {
	//logger log.Log
}

//
// This file contains the true implementations--the server side--for the methods
// defined in business.proto.
//

func (m *myServer) MediaTypesInStock(pctx lib.Pctx) (proto.Message, error) {
	return &pb.MediaTypesInStockResponse{
		InStock: []pb.Media{pb.Media_MEDIA_BETA, pb.Media_MEDIA_CASSETTE, pb.Media_MEDIA_LASER_DISC, pb.Media_MEDIA_VHS, pb.Media_MEDIA_VINYL, pb.Media_MEDIA_ATARI_CART},
	}, nil
}

func (m *myServer) BestOfAllTime(pctx lib.Pctx, inProto proto.Message) (proto.Message, error) {
	in := inProto.(*pb.BestOfAllTimeRequest)

	out := &pb.BestOfAllTimeResponse{
		Item: &pb.Item{},
	}
	m.log(pctx, pblog.LogLevel_LOGLEVEL_DEBUG, "reached BestOfAllTime, computing his choices")

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
	m.log(pctx, pblog.LogLevel_LOGLEVEL_INFO, "unexpected content type in request %d", int32(in.Ctype))
	return nil, fmt.Errorf("unexpected content type request int %d", int32(in.Ctype))
}

func (m *myServer) Revenue(pctx lib.Pctx, in proto.Message) (proto.Message, error) {
	out := &pb.RevenueResponse{}
	m.log(pctx, pblog.LogLevel_LOGLEVEL_WARNING, "Revenue() not yet implemented, ignoring input value, returning dummy values")
	out.Revenue = 817.71
	return out, nil
}

func (m *myServer) SoldItem(pctx lib.Pctx, in proto.Message) error {
	m.log(pctx, pblog.LogLevel_LOGLEVEL_WARNING, "SoldItem() not yet implemented, ignoring input value")
	return nil
}

// Ready is a check, if this returns false the library should abort and not attempt to run this service.
// Normally this is used to inform the kernel that we are exporting some package and that we are ready to
// run.
func (m *myServer) Ready() bool {

	// logger, err := log.LocateLog()
	// if err != nil {
	// 	print("ERROR trying to create log client: ", err.Error(), "\n")
	// 	return false
	// }
	// m.logger = logger

	if _, err := lib.Export1("demo.vvv", "Store"); err != nil {
		print("ready: error in attempt to export demo.vvv: ", err.Error(), "\n")
		return false
	}
	if _, err := lib.Run(false); err != nil {
		print("ready: error in attempt to signal Run: ", err.Error(), "\n")
		return false
	}
	return true

}

func (m *myServer) log(pctx lib.Pctx, level pblog.LogLevel, spec string, rest ...interface{}) {
	// msg := fmt.Sprintf(spec, rest...)
	// req := pblog.LogRequest{
	// 	Stamp:   timestamppb.New(pctx.Now()),
	// 	Level:   level,
	// 	Message: msg,
	// }
	print("log request in server\n")
	// m.logger.Log(&req)
}

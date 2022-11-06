package main

import (
	"fmt"

	"demo/vvv/proto/g/vvv"
	"demo/vvv/proto/g/vvv/pb"

	"github.com/iansmith/parigot/g/pb/kernel"
	"github.com/iansmith/parigot/g/pb/log"
	"github.com/iansmith/parigot/lib"
	"google.golang.org/protobuf/proto"
)

func main() {
	vvv.Run(&myServer{})
}

// this type better implement vvv.StoreServer
type myServer struct {
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
	pctx.Log(log.LogLevel_LOGLEVEL_DEBUG, "reached BestOfAllTime, computing his choices")

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
	pctx.Log(log.LogLevel_LOGLEVEL_INFO, fmt.Sprintf("unexpected content type in request %d", int32(in.Ctype)))
	return nil, fmt.Errorf("unexpected content type request int %d", int32(in.Ctype))
}

func (m *myServer) Revenue(pctx lib.Pctx, in proto.Message) (proto.Message, error) {
	out := &pb.RevenueResponse{}
	pctx.Log(log.LogLevel_LOGLEVEL_WARNING, "Revenue() not yet implemented, ignoring input value, returning dummy values")
	out.Revenue = 817.71
	return out, nil
}

func (m *myServer) SoldItem(pctx lib.Pctx, in proto.Message) error {
	pctx.Log(log.LogLevel_LOGLEVEL_WARNING, "SoldItem() not yet implemented, ignoring input value")
	return nil
}

// Ready is a check, if this returns false the library should abort and not attempt to run this service.
// Normally this is used to inform the kernel that we are exporting some package and that we are ready to
// run.
func (m *myServer) Ready() bool {

	if _, err := lib.Export1("demo.vvv", "Store"); err != nil {
		return false
	}
	if _, err := lib.Run(&kernel.RunRequest{Wait: false}); err != nil {
		return false
	}
	return true

}

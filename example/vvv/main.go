package main

import (
	"fmt"

	"demo/vvv/proto/g/vvv/pb"

	"github.com/iansmith/parigot/g/pb/log"
	"github.com/iansmith/parigot/lib"
	"google.golang.org/protobuf/proto"
)

func main() {
	run(&myServer{})
}

// this type better implement StoreServer
type myServer struct {
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

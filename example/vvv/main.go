package main

import (
	"demo/vvv/proto/gen/demo/vvv"
)

// server side

func main() {
	// should be generated, just like locate
	vinnysHandler := vvv.StoreHandler{
		BestOfAllTime: BestOfAllTime,
	}
	// should be generated, just like locate
	vvv.RegisterStoreHandler(vinnysHandler)
}
func BestOfAllTime() *vvv.Media {
	return &vvv.Media{
		Title: "The Queen Is Dead",
	}
}

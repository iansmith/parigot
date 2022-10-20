package main

import (
	"proto/g/demo/vvv"
)

// server side

func main() {
	// should be generated, just like locate
	vinnysHandler := vvv.VinnysStoreHandler{
		BestOfAllTime: BestOfAllTime,
	}
	// should be generated, just like locate
	vvv.RegisterVinnysStoreHandler(vinnysHandler)
}
func BestOfAllTime() *vvv.Media {
	return &vvv.Media{
		Title: "The Queen Is Dead",
	}
}

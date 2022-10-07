package main

import (
	greet "example/greet/proto/gen/greet"
	"log"
)

func main() {
	var g
	g = lib_Connect[Greet]("demo", "greet")
	result, err := g.Greet(greet.GreetRequest{Name: "fleazil"})
	log.Println(result.Msg.Greeting)
}

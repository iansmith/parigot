package main

import (
	"context"
	"log"
	"net/http"

	greetv1 "example/greet/proto/gen/greet"
	greetv1connect "example/greet/proto/gen/greet/greetconnect"

	"github.com/bufbuild/connect-go"
)

func main() {
	client := greetv1connect.NewGreetClient(
		http.DefaultClient,
		"http://localhost:8080",
		// connect.WithGRPC(), if you uncomment this, you can use GRPC
	)
	res, err := client.Greet(
		context.Background(),
		connect.NewRequest(&greetv1.GreetRequest{Name: "Jane"}),
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Msg.Greeting)
}

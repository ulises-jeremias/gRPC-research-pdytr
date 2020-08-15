package main

import (
	"../proto"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:4444", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	c := proto.NewOperationsClient(conn)
	req := &proto.Request{A: int64(25), B: int64(10)}

	suma, err := c.Add(context.Background(), req)
	if err != nil {
		log.Fatalf("Error when calling Add: %s", err)
	}
	log.Printf("Response from server: %d", int64(suma.Result))

	resta, err := c.Sub(context.Background(), req)
	if err != nil {
		log.Fatalf("Error when calling Sub: %s", err)
	}
	log.Printf("Response from server: %d", int64(resta.Result))

}
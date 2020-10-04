package main

import (
	"log"

	simple ".."

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:4444", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	c := simple.NewOperationsClient(conn)
	req := &simple.Request{A: int64(25), B: int64(10)}

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

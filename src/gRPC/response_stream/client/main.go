package main

import (
	"fmt"
	"io"
	"log"

	response_stream ".."

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:4444", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := response_stream.NewOperationsClient(conn)

	//Requesting data on screen
	var num int64
	fmt.Println("Enter a number and the server will send the first 10 multiples of the same")
	fmt.Scanf("%d\n", &num)

	//Create a stream
	req := &response_stream.Request{Num: num}
	stream, err := client.Multiple(context.Background(), req)
	if err != nil {
		log.Fatalf("Open stream error %v", err)
	}

	//Receiving data from server
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Printf("%d -> ", reply.OneMultiple)
	}
}

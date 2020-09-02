package main

import (
	".."
	"log"
	"math/rand"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:4444", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := request_stream.NewOperationsClient(conn)

	// Create a stream
	stream, err := client.Summation(context.Background())
	if err != nil {
		log.Fatalf("openn stream error %v", err)
	}

	//sending 10 random numbers
	var local_result int64 = 0
	for i := 1; i <= 10; i++ {
		// generate random nummber and send it to stream
		rnd := int64(rand.Intn(i))
		req := request_stream.Request{Num: rnd}
		if err := stream.Send(&req); err != nil {
			log.Fatalf("can not send %v", err)
		}
		log.Printf("sending number --> %d ", rnd)
		//performing a local sum to verify the server results
		local_result += req.Num
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Printf("sum of the server: %d , and local sum: %d", reply.Result, local_result)

	}
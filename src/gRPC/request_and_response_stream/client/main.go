package main

import (
	"fmt"
	"io"
	"log"

	request_and_response_stream ".."

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:4444", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := request_and_response_stream.NewOperationsClient(conn)

	// Create a stream
	stream, err := client.Max(context.Background())
	done := make(chan bool)
	if err != nil {
		log.Fatalf("openn stream error %v", err)
	}

	// sending numbers
	var num int64
	var ok bool = true
	fmt.Prinln("Starting a bidirection conection with the server")
	fmt.Prinln("Enter numbers, the server will compare the numbers that reach you with your local maximum, if the number sent is greater than the local maximum, then it will be updated and the client will be notified")
	for ok {
		// Requesting data to the user
		fmt.Prinln("Enter a number (0 to finish)")
		fmt.Scanf("%d \n", &num)

		if num == 0 {
			ok = false
			//closing sending stream
			stream.CloseSend()
			fmt.Prinln("Conection closed")
			continue
		}

		// send request data
		req := request_and_response_stream.Request{Num: num}
		if err := stream.Send(&req); err != nil {
			log.Fatalf("can not send %v", err)
		}

		// receive res data
		reply, err := stream.Recv()
		if err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Printf("Actual max is: %d \n", reply.Result)

	}

	// closing receive stream
	_, err := stream.Recv()
	if err == io.EOF {
		close(done)
		return
	}
	if err != nil {
		log.Fatalf("Can not receive %v", err)
	}
}

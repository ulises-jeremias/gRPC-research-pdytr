package main

import (
	".."
	"log"
	"fmt"
	"io"

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
	fmt.Printf("Ingrese un numero y el servidor enviara los primeros 10 multiplos del mismo \n")
	fmt.Scanf("%d \n", &num)

	//Create a stream
	req := &response_stream.Request{Num: num}
	stream, err := client.Multiple(context.Background(), req)
	if err != nil {
		log.Fatalf("openn stream error %v", err)
	}

	//Receiving data from server
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
		}
		fmt.Printf("%d -> ",reply.OneMultiple)
	}
	}

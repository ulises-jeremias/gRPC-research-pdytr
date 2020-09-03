package main

import (
	".."
	"log"
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:5555", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := request_stream.NewOperationsClient(conn)

	// Create a stream
	stream, err := client.Summation(context.Background())
	if err != nil {
		log.Fatalf("openn stream error %v", err)
	}

	//sending 10 numbers
	var num int64
	for i := 1; i <= 10; i++ {
		//Requesting data on screen
		fmt.Printf("Ingrese un numero\n")
		fmt.Scanf("%d \n", &num)
		req := request_stream.Request{Num: num}
		if err := stream.Send(&req); err != nil {
			log.Fatalf("can not send %v", err)
		}
		fmt.Printf("Faltan enviar %d numeros mas al servidor \n", (10-i))
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Printf("Respuesta del servidor: \n Suma de los 10 numeros: %d ", reply.Result)

	}
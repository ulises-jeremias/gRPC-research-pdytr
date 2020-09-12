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

	client := request_and_response_stream.NewOperationsClient(conn)

	// Create a stream
	stream, err := client.Max(context.Background())
	done := make(chan bool)
	if err != nil {
		log.Fatalf("openn stream error %v", err)
	}

	//sending numbers
	var num int64
	var ok bool = true
	fmt.Printf("Empezando una transmisión bidireccional con el servidor......\n\n")
	fmt.Printf("Ingrese numeros, el servidor ira comparando los numeros que le lleguen con su maximo local, si el numero enviado es mayor al maximo local, entonces se actualizara el mismo y se notificara al cliente\n\n")
	for ok {
		//Requesting data on screen
		fmt.Printf("Ingrese un numero. (0 para finalizar)\n\n")
		fmt.Scanf("%d \n", &num)

		if num == 0{
			ok = false
			//closing sending stream
			stream.CloseSend()
			fmt.Printf("Conexion terminada \n")
			continue
		}

		//send
		req := request_and_response_stream.Request{Num: num}
		if err := stream.Send(&req); err != nil {
			log.Fatalf("can not send %v", err)
		}

		//recive
		reply, err := stream.Recv()
		if err != nil {
			log.Fatalf("%v",err)
		}
		fmt.Printf("El maximo actual del servidor es: %d \n",reply.Result)
		
	}
	//closing receive stream
	reply, err := stream.Recv()
	_ = reply 
			if err == io.EOF {
				close(done)
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
}
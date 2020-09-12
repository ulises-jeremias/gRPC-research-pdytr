package main

import (
	".."
	"net"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func main() {
	listener, err := net.Listen("tcp", ":4444")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	request_and_response_stream.RegisterOperationsServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}

func (s *server) Max(stream request_and_response_stream.Operations_MaxServer) error {

	log.Println("start new server")
	var max int64

	for {
		// receive data from stream
		req, err := stream.Recv()
		if err == io.EOF {
			// return will close stream from server side
			log.Println("exit")
			return nil
		}
		if err != nil {
			log.Printf("receive error %v", err)
			return err
		}
		
		// continue if number reveived from stream
		// less than max
		if req.Num > max {
			max = req.Num
		}

		//send to stream
		resp := request_and_response_stream.Response{Result: max}
		if err := stream.Send(&resp); err != nil {
			log.Printf("send error %v", err)
		}
		log.Printf("send new max=%d", max)
	}
}

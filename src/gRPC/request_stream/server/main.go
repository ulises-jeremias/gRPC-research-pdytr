package main

import (
	"io"
	"net"

	request_stream ".."

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func main() {
	listener, err := net.Listen("tcp", ":5555")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	request_stream.RegisterOperationsServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}

func (s *server) Summation(stream request_stream.Operations_SummationServer) error {
	var result int64 = 0
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&request_stream.Response{Result: result})
		}
		if err != nil {
			return err
		}
		result += request.Num
	}
}

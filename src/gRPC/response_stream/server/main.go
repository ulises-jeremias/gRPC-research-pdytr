package main

import (
	"net"

	response_stream ".."

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
	response_stream.RegisterOperationsServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}

func (s *server) Multiple(request *response_stream.Request, stream response_stream.Operations_MultipleServer) error {
	var aux int64 = 0
	for i := 1; i <= 10; i++ {
		aux += request.Num
		resp := response_stream.Response{OneMultiple: aux}
		if err := stream.Send(&resp); err != nil {
			return err
		}
	}
	return nil
}

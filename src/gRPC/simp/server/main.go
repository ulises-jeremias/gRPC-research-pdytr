package main

import (
	".."
	"context"
	"net"

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
	simp.RegisterOperationsServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}

func (s *server) Add(ctx context.Context, request *simp.Request) (*simp.Response, error) {
	a, b := request.GetA(), request.GetB()

	result := a + b

	return &simp.Response{Result: result}, nil
}

func (s *server) Sub(ctx context.Context, request *simp.Request) (*simp.Response, error) {
	a, b := request.GetA(), request.GetB()

	result := a - b

	return &simp.Response{Result: result}, nil
}

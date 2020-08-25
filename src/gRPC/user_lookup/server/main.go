package main

import (
	".."
	"context"
	"net"
	"os/user"

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
	user_lookup.RegisterOperationsServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}

func (s *server) Byname(ctx context.Context, username *user_lookup.Username) (*user_lookup.UserId, error) {
	name := username.GetName()

	user_found, err  := user.Lookup(name)
	if err != nil {
		panic(err)
	}
	
	return &user_lookup.UserId{Num: user_found.Uid}, nil
}

func (s *server) Bynum(ctx context.Context, user_id *user_lookup.UserId) (*user_lookup.Username, error) {
	num := user_id.GetNum()

	user_found, err := user.LookupId(num)
	if err != nil {
		panic(err)
	}

	return &user_lookup.Username{Name: user_found.Username}, nil
}

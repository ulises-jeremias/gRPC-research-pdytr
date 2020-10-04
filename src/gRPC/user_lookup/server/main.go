package main

import (
	"context"
	"net"
	"os/user"

	user_lookup ".."

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

func (s *server) ByUsername(ctx context.Context, username *user_lookup.Username) (*user_lookup.UserID, error) {
	name := username.GetName()

	userFound, err := user.Lookup(name)
	if err != nil {
		panic(err)
	}

	return &user_lookup.UserID{Num: userFound.Uid}, nil
}

func (s *server) ByID(ctx context.Context, userID *user_lookup.UserID) (*user_lookup.Username, error) {
	num := userID.GetNum()

	userFound, err := user.LookupId(num)
	if err != nil {
		panic(err)
	}

	return &user_lookup.Username{Name: userFound.Username}, nil
}

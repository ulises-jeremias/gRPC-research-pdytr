package main

import (
	"log"

	user_lookup ".."

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:4444", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	c := user_lookup.NewOperationsClient(conn)
	username := &user_lookup.Username{Name: "root"}

	userID, err := c.ByUsername(context.Background(), username)
	if err != nil {
		log.Fatalf("Error when calling ByUsername: %s", err)
	}
	log.Printf("Root user ID: %s", userID.Num)

	userID = &user_lookup.UserID{Num: "0"}
	username, err = c.ByID(context.Background(), userID)
	if err != nil {
		log.Fatalf("Error when calling Sub: %s", err)
	}
	log.Printf("Username of user with ID 0: %s", username.Name)
}

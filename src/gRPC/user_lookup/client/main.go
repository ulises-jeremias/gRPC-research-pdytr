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

	var userID, err = c.ByUsername(context.Background(), username)
	if err != nil {
		log.Fatalf("Error when calling ByUsername: %s", err)
	}
	log.Printf("El ID de root es: %s", userID.Num)

	userID = &user_lookup.UserID{Num: "0"}
	username2, err := c.ById(context.Background(), userID)
	if err != nil {
		log.Fatalf("Error when calling Sub: %s", err)
	}
	log.Printf("El username del Id 0 es: %s", username2.Name)

}

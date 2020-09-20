package main

import (
	".."
	"log"

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

	user_id, err := c.ByUsername(context.Background(), username)
	if err != nil {
		log.Fatalf("Error when calling ByUsername: %s", err)
	}
	log.Printf("El ID de root es: %s", user_id.Num)

	//nose por que no me deja reutilizar las varialbes, no te enojes uli :(
	user_id2 := &user_lookup.UserId{Num: "0"}
	username2, err := c.ById(context.Background(), user_id2)
	if err != nil {
		log.Fatalf("Error when calling Sub: %s", err)
	}
	log.Printf("El username del Id 0 es: %s", username2.Name)

}

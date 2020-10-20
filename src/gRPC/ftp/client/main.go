package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	ftp ".."

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type struct Command {
	name string
	description string
	handle func(args HandlerArgs)
}

type struct HandlerArgs {
	verbose bool
	bytes int64
	initial_pos int64
	src string
	dest string
}

func main() {
	conn, err := grpc.Dial("localhost:4444", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := request_and_response_stream.NewOperationsClient(conn)

	verbose := flag.Bool("verbose", false, "")
	src := flag.String("src", "", "")
	dest := flag.String("dest", "", "")
	initial_pos := flag.Int64("init-pos", 0, "")
	bytes := flag.Int64("bytes", "-1", "")
	list := flag.Bool("list", false, "")

	flag.Parse()

	commands := map[string]HandlerArgs{
		"write": {
			name = "write",
			description = "Add a file from --src to --dest",
			handle = ftp_write,
		},
		"read": {
			name = "read",
			description = "Store a file from --src to --dest",
			handle = ftp_read,
		},
		"list": {
			name = "list",
			description = "List all files from --src",
			handle = ftp_list,
		}
	}

	if len(os.Args) < 2 {
		log.Printf("Usage: %s\n", argv[0])
		for k, c := range commands {
				log.Printf("\t- %s: %s\n", c.name, c.description)
		}
		os.Exit()
	}

	command_name := os.Args[1]
	if command, ok := commands[command]; ok {
			
	}
}

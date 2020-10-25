package main

import (
	"flag"
	"fmt"
	"os"

	ftp ".."
	command "./command"

	"google.golang.org/grpc"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [flags] OPERATION\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOperations:\n\n")
	for _, c := range command.Commands {
		fmt.Fprintf(os.Stderr, "  %s: %s\n", c.Name, c.Description)
	}
	fmt.Fprintf(os.Stderr, "\nFlags:\n\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	verbose := flag.Bool("v", false, "")
	src := flag.String("src", "", "")
	dest := flag.String("dest", "tmp1", "")
	initialPos := flag.Int64("i", 0, "")
	bytes := flag.Int64("b", -1, "")
	list := flag.Bool("l", false, "")

	flag.Usage = usage
	flag.Parse()

	if len(os.Args) < 2 {
		usage()
	}

	conn, err := grpc.Dial("localhost:4444", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := ftp.NewOperationsClient(conn)

	commandName := os.Args[len(os.Args)-1]
	if c, ok := command.Commands[commandName]; ok {
		args := command.HandlerArgs{
			Verbose:    *verbose,
			Bytes:      *bytes,
			List:       *list,
			InitialPos: *initialPos,
			Src:        *src,
			Dest:       *dest,
		}

		c.Handle(client, args)
	}
}

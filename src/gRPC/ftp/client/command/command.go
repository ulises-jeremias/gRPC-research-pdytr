package command

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	ftp "../.."
)

// Command type that contains name, description and handler type
type Command struct {
	Name        string
	Description string
	Handle      func(client ftp.OperationsClient, args HandlerArgs)
}

// HandlerArgs type that contains needed data to be passed to a command
type HandlerArgs struct {
	Verbose    bool
	List       bool
	Bytes      int64
	InitialPos int64
	Src        string
	Dest       string
}

// Commands map struct containing commands
var Commands = map[string]Command{
	"write": {
		Name:        "write",
		Description: "Add a file from --src to --dest",
		Handle:      ftpWrite,
	},
	"read": {
		Name:        "read",
		Description: "Store a file from --src to --dest",
		Handle:      ftpRead,
	},
	"list": {
		Name:        "list",
		Description: "List all files",
		Handle:      ftpList,
	},
}

func ftpRead(client ftp.OperationsClient, args HandlerArgs) {
	if args.Verbose {
		log.Println(args)
	}

	pos := args.InitialPos
	readAll := args.Bytes == -1
	bytes := args.Bytes

	// send request
	req := ftp.ReadRequest{
		Name:  args.Src,
		Pos:   pos,
		Bytes: bytes,
	}

	// Create a stream channel
	stream, err := client.Read(context.Background(), &req)
	done := make(chan bool)
	if err != nil {
		log.Fatalf("open stream error %v", err)
	}

	filePath := path.Join(args.Dest)
	dirPath := filepath.Dir(filePath)

	// ensure that dir exists with perm os.ModePerm: 0777
	_ = os.MkdirAll(dirPath, os.ModePerm)

	// open file
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalf("Error to open [file=%v]: %v", filePath, err.Error())
	}

	for {
		if !readAll && bytes <= 0 {
			break
		}

		// receive data
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("can not receive %v", err)
		}

		// save buffer data
		dataLen := len(res.Data)
		buf := []byte(res.Data)

		// write file content using buffer
		n, err := f.Write(buf)
		if n == 0 {
			if err == nil {
				// ignore
			} else {
				log.Printf("Error %v", err)
				break
			}
		}

		pos += int64(dataLen)
		if !readAll {
			bytes -= int64(dataLen)
		}

		if !res.ContinueReading {
			break
		}
	}

	defer f.Close()
	close(done)
}

func ftpWrite(client ftp.OperationsClient, args HandlerArgs) {
	// Create a stream channel
	stream, err := client.Write(context.Background())
	done := make(chan bool)
	if err != nil {
		log.Fatalf("open stream error %v", err)
	}

	if args.Verbose {
		log.Println(args)
	}

	f, err := os.Open(args.Src)
	if err != nil {
		log.Fatalf("Error to read [file=%v]: %v", args.Src, err.Error())
	}
	if _, err := f.Seek(args.InitialPos, 0); err != nil {
		log.Fatalf("Error to seek position %d for file [file=%v]: %v", args.InitialPos, args.Src, err.Error())
	}

	totalBytesReaded := int64(0)

	for {
		if args.Bytes != -1 && totalBytesReaded >= args.Bytes {
			break
		}

		r := bufio.NewReader(f)
		buf := make([]byte, 4*1024)

		// read file content using buffer
		n, err := r.Read(buf)
		buf = buf[:n]
		if n == 0 {
			log.Print(totalBytesReaded)
			if err == nil {
				break
			}
			if err == io.EOF {
				break
			}
			log.Fatalf("Error %v", err)
		}

		defer f.Close()

		// save buffer data
		dataVal := string(buf)

		totalBytesReaded += int64(n)

		// send request
		req := ftp.WriteRequest{
			Name:     args.Dest,
			Data:     dataVal,
			Checksum: ftp.Djb2(buf),
		}
		if err := stream.Send(&req); err != nil {
			log.Fatalf("can not send %v", err)
		}
	}

	// closing receive stream
	_, err = stream.CloseAndRecv()
	if err == io.EOF {
		close(done)
		return
	}
	if err != nil {
		log.Fatalf("can not receive %v", err)
	}
}

func ftpList(client ftp.OperationsClient, args HandlerArgs) {
	req := &ftp.ListRequest{Name: args.Src, List: args.List}

	res, err := client.List(context.Background(), req)
	if err != nil {
		log.Fatalf("Error when calling Add: %s", err)
	}
	log.Printf("%s", res.Paths)
}

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
	// Create a stream channel
	stream, err := client.Read(context.Background())
	done := make(chan bool)
	if err != nil {
		log.Fatalf("open stream error %v", err)
	}

	if args.Verbose {
		log.Println(args)
	}

	filePath := path.Join(args.Dest)
	dirPath := filepath.Dir(filePath)

	// ensure that dir exists with perm os.ModePerm: 0777
	_ = os.MkdirAll(dirPath, os.ModePerm)

	// open file
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalf("Error to open [file=%v]: %v", filePath, err.Error())
	}

	pos := args.InitialPos
	readAll := args.Bytes == -1
	bytes := args.Bytes

	for {
		if !readAll && bytes <= 0 {
			break
		}

		// send request
		req := ftp.ReadRequest{
			Name:  args.Src,
			Pos:   pos,
			Bytes: bytes,
		}
		if err := stream.Send(&req); err != nil {
			log.Fatalf("can not send %v", err)
		}

		// receive data
		res, err := stream.Recv()
		if err != nil {
			log.Fatalf("%v", err)
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

	// closing receive stream
	_, err = stream.Recv()
	if err == io.EOF {
		close(done)
		return
	}
	if err != nil {
		log.Fatalf("can not receive %v", err)
	}
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
		buf := make([]byte, 0, 1024)

		// read file content using buffer
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
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

		writeMode := os.O_WRONLY
		if totalBytesReaded > 0 {
			writeMode = os.O_APPEND
		}

		// send request
		req := ftp.WriteRequest{
			Name:     args.Dest,
			Data:     dataVal,
			Mode:     int32(writeMode),
			Checksum: ftp.Djb2(buf),
		}
		if err := stream.Send(&req); err != nil {
			log.Fatalf("can not send %v", err)
		}

		// receive data
		_, err = stream.Recv()
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	// closing receive stream
	_, err = stream.Recv()
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

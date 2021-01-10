package main

import (
	"context"

	ftp ".."

	"bufio"
	"io"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"

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
	ftp.RegisterOperationsServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}

func (s *server) Read(req *ftp.ReadRequest, stream ftp.Operations_ReadServer) error {
	isEOF := false
	bytes := req.Bytes
	filePath := path.Join("store", req.Name)

	log.Printf("Path: %s", filePath)

	// limit bytes value to save on chunck
	// with size under 1025
	if bytes > 4*1024 || bytes < 0 {
		bytes = 4 * 1024
	}

	// open file
	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error to read [file=%v]: %v", filePath, err.Error())
		return err
	}
	if _, err := f.Seek(req.Pos, 0); err != nil {
		log.Printf("Error to seek position %d for file [file=%v]: %v", req.Pos, filePath, err.Error())
		return err
	}

	for {

		r := bufio.NewReader(f)
		buf := make([]byte, bytes)

		// read file content using buffer
		n, err := r.Read(buf)
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				break
			} else if err == io.EOF {
				isEOF = true
				break
			} else {
				log.Printf("Error %v", err)
				return err
			}
		}

		// save buffer data
		dataVal := string(buf)
		dataLen := int64(len(buf))

		// send response to stream
		res := ftp.ReadResponse{
			Data:            dataVal,
			Name:            req.Name,
			ContinueReading: !isEOF,
		}
		if err := stream.Send(&res); err != nil {
			log.Printf("Error %v", err)
		}
		log.Printf("Sent chunk of size = %d for file %s", dataLen, req.Name)
	}

	defer f.Close()
	return nil
}

func (s *server) Write(stream ftp.Operations_WriteServer) error {

	for {
		// receive data from stream
		req, err := stream.Recv()
		if err == io.EOF {
			// send response to stream
			res := ftp.WriteResponse{}
			if err := stream.SendAndClose(&res); err != nil {
				log.Printf("Error %v", err)
			}
			return nil
		}
		if err != nil {
			log.Printf("Error: Receive error %v", err)
			return err
		}

		filePath := path.Join("store", req.Name)
		dirPath := filepath.Dir(filePath)

		// save buffer data
		dataLen := len(req.Data)
		buf := []byte(req.Data)
		buf = buf[:dataLen]

		log.Printf("Path: %s", filePath)

		// ensure that dir exists with perm os.ModePerm: 0777
		_ = os.MkdirAll(dirPath, os.ModePerm)

		// open file
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			log.Printf("Error to open [file=%v]: %v", filePath, err.Error())
			return err
		}

		// check checksum
		checksum := ftp.Djb2(buf)
		if req.Checksum != checksum {
			log.Printf("Error in checksum!!\nOriginal: %d\nOwn: %d\n", req.Checksum, checksum)
			return err
		}

		// write file content using buffer
		n, err := f.Write(buf)
		if n == 0 {
			if err == nil {
				// ignore
			}
			log.Printf("Error %v", err)
			return err
		}

		defer f.Close()
		log.Printf("Wrote chunk of size = %d for file %s", dataLen, req.Name)
	}
}

func (s *server) List(ctx context.Context, req *ftp.ListRequest) (res *ftp.ListResponse, err error) {

	var files []string

	root := "store"
	err = filepath.Walk(root, ftp.FileHandler(&files))
	if err != nil {
		log.Printf("Error %v", err)
		return nil, err
	}

	sep := "\t"
	if req.List {
		sep = "\n"
	}

	// send response to stream
	res = &ftp.ListResponse{
		Paths: strings.Join(files, sep) + "\n",
	}
	return res, nil
}

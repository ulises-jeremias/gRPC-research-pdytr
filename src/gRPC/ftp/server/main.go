package main

import (
	"fmt"

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
	ftp_stream.RegisterOperationsServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}

func (s *server) Read(req *ftp_stream.ReadRequest, stream ftp_stream.Operations_ReadServer) error {

	for {
		// receive data from stream
		req, err := stream.Recv()
		if err == io.EOF {
			// return will close stream from server side
			log.Println("exit")
			return nil
		}
		if err != nil {
			log.Printf("Error: Receive error %v", err)
			return err
		}

		isEOF := false
		bytes := req.Bytes
		filePath := path.Join("store", req.Name)

		log.Println("Path: %s", filePath)

		// limit bytes value to save on chunck
		// with size under 1025
		if bytes > 1024 {
			bytes = 1024
		}

		// open file
		f, err := os.Open(filePath)
		if err != nil {
			log.Println("Error to read [file=%v]: %v", filePath, err.Error())
			return err
		}
		if _, err := input.Seek(req.Pos, 0); err != nil {
			log.Println("Error to seek position %d for file [file=%v]: %v", req.Pos, filePath, err.Error())
			return err
		}

		r := bufio.NewReader(f)
		buf := make([]byte, 0, bytes)

		// read file content using buffer
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				// ignore
			}
			if err == io.EOF {
				isEOF = true
			}
			log.Println("Error %v", err)
			return err
		}

		// save buffer data
		dataVal := string(buf)
		dataLen := int64(len(buf))

		// send response to stream
		resp := ftp_stream.ReadResponse{
			Data:            dataVal,
			Name:            name,
			ContinueReading: !isEOF,
		}
		if err := stream.Send(&resp); err != nil {
			log.Printf("Error %v", err)
		}
		log.Printf("Sent chunk of size = %d for file %s", dataLen, name)
	}
}

func (s *server) Write(req *ftp_stream.WriteRequest, stream ftp_stream.Operations_WriteServer) error {

	for {
		// receive data from stream
		req, err := stream.Recv()
		if err == io.EOF {
			// return will close stream from server side
			log.Println("exit")
			return nil
		}
		if err != nil {
			log.Printf("Error: Receive error %v", err)
			return err
		}

		filePath := path.Join("store", req.Name)
		dirPath := filepath.Dir(filePath)

		// save buffer data
		dataLen := len(req.data)
		buf := make([]byte, 0, 1026)
		copy(buf[:], req.data)

		log.Println("Path: %s", filePath)

		// ensure that dir exists with perm os.ModePerm: 0777
		_ := os.MkdirAll(dirPath, os.ModePerm)

		// open file
		f, err := os.Open(filePath, req.Mode)
		if err != nil {
			log.Println("Error to open [file=%v]: %v", filePath, err.Error())
			return err
		}

		// check checksum
		checksum := ftp.Djb2(buf)
		if incomeChecksum != checksum {
			log.Println("Error in checksum!!\nOriginal: %d\nOwn: %d\n", req.Checksum, checksum)
			return err
		}

		// write file content using buffer
		w := bufio.NewWriter(f)
		n, err := w.Write(buf[:cap(buf)])
		if n == 0 {
			if err == nil {
				// ignore
			}
			log.Println("Error %v", err)
			return err
		}

		// send response to stream
		resp := ftp_stream.WriteResponse{}
		if err := stream.Send(&resp); err != nil {
			log.Printf("Error %v", err)
		}
		log.Printf("Wrote chunk of size = %d for file %s", dataLen, name)
	}
}

func fileHandler(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		*files = append(*files, path.Name())
		return nil
	}
}

func (s *server) List(req *ftp_stream.ListRequest, stream ftp_stream.Operations_WriteServer) error {
	var files []string

	root := "store"
	err := filepath.Walk(root, fileHandler(&files))
	if err != nil {
		panic(err)
	}

	sep := '\t'
	if req.list {
		sep = '\n'
	}

	// send response to stream
	resp := ftp_stream.ListResponse{
		paths: strings.Join(files, sep) + '\n'
	}
	if err := stream.Send(&resp); err != nil {
		log.Printf("Error %v", err)
	}
}

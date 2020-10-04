# gRPC-research-pdytr

## Quickstart

We provide a docker environment to execute the different applications. To start the docker container execute the following command,

```sh
$ ./bin/grpc --app=<app> [--build] [--root]
```

where,

```
app = simp | ftp
```

_NOTE: Execute `./bin/grpc -h` to know more about flags._

This will execute a docker container in interactive mode with a volume in the directory of the application.

<p align="center">
  <img src="./static/terminal.png">
</p>

Once the docker container is running you will be able to run the following examples!

## Simple example

```sh
# generate proto file
~> protoc --go_out=plugins=grpc:. ./simp.proto

# execute server in background
~> go run ./server/main.go &

# execute client
~> go run ./client/main.go
```

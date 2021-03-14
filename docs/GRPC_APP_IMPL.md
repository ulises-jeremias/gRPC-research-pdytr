## Implementación de Cliente/Servidor FTP con gRPC

Este es un tutorial básico de introducción a gRPC en Go.

Esta sección propociona una guía básica de programación en GO con gRPC.

Al seguir esta guía se puede aprender a:

- Defina un servicio en un archivo `.proto`.
- Genere código de cliente y servidor utilizando el compilador `protoc`.
- Utilice la API de GO gRPC para escribir un cliente y un servidor sencillos para su servicio.

Para esto se asume que ya se ha realizado alguna de las opciones de setup de GRPC descriptas en la [guía de instalación](./GRPC_APP_RUN.md).

En las siguientes subsecciones vamos a repasar el diseño e implementación de un cliente y servidor FTP escrito en GO con gRPC. La implementación final puede encontrarse [aquí](https://github.com/ulises-jeremias/gRPC-research-pdytr/tree/master/src/gRPC/ftp).

Para esto creamos el directorio `./src/gRPC/ftp` lo cual nos permite levantar el ambiente docker de forma rápida utilizando el siguiente comando:

```sh
$ ./bin/grpc --app=ftp --root
```

o simplemente ejecutar todo localmente en caso de no querer utilizar docker.

### Diseño

Esta implementación es una versión resumida del protocolo FTP en la cual se utilizan cada uno de los conceptos con los cuales se experimenta en los ejemplos anteriores.

Los métodos implementados son **read**, **write** y **list**.

#### Read

La operación READ se implementa de la siguiente forma:

1. El cliente solicita la lectura de un chunck de _b bytes_ a partir de una posición **pos** en un archivo de nombre **name**.

2. En caso de que exista el recurso, el servidor devuelve el chunck solicitado, **data**.

3. En cliente continua recibiendo porciones del archivo hasta que el servidor responde con **continue_reading** igual a _false_, o retorne algún error.

Como se puede notar de esta secuencia, la misma resulta en una interacción en la cual el cliente solicita el recurso una única vez al inicio y el servidor envía porciones del mismo hasta que todos los chuncks sean leidos correctamente o exista algún error. Teniendo en cuenta que la cantidad de interacciones varía según el tamaño del archivo y tamaño de cada chunck es que pensamos que sería interesante y óptimo utilizar algún modelo de streaming, y particularmente para este caso, utilizar el modelo **Server Side Streaming** para la definición de esta operación.

Para todo esto definimos dos tipos de mensajes: `ReadRequest` y `ReadResponse` con los datos necesarios, y el servicio `Read` de la siguiente forma:

```proto
syntax = "proto3";

package ftp;

message ReadRequest {
    string name = 1;
    int64 pos = 2;
    int64 bytes = 3;
}

message ReadResponse {
    string name = 1;
    string data = 2;
    bool continue_reading = 3;
}

service Operations {
    rpc Read(ReadRequest) returns (stream ReadResponse);
}
```

#### Write

La operación WRITE se implementa de la siguiente forma:

1. El cliente envía un chunck _b bytes_ a ser almacenado a partir de una posición **pos** en un archivo de nombre **name**. A su vez se envía un checksum para validar el contenido enviado.

2. En caso de poder utilizar el recurso asignado, el servidor almacena el chunck enviado, **data**.
3. El cliente continua enviando porciones del archivo hasta que se haya enviado todo el contenido del archivo.
4. El cliente sigue enviando porciones de archivos. Al finalizar espera la respuesta del servidor para comprobar si existió algún error.

A deferencia del caso anterior, la operación se define como un proceso en el cual el cliente envia de forma continua datos a escribir en el servidor y recibe un resumen del total de operaciones en caso de haber un error. Teniendo en cuenta que la cantidad de interacciones varía según el tamaño del archivo y tamaño de cada chunck es que pensamos que sería interesante y óptimo utilizar un modelo de stream, pero a diferencia del caso anterior, este será un modelo de **Client Side Streaming**.

Para todo esto definimos dos tipos de mensajes: `WriteRequest` y `WriteResponse` con los datos necesarios, y el servicio `Write` de la siguiente forma:

```proto
syntax = "proto3";

package ftp;

message WriteRequest {
    string name = 1;
    string data = 2;
    int64 checksum = 3;
}

message WriteResponse {}

service Operations {
    rpc Write(stream WriteRequest) returns (WriteResponse);
}
```

#### List

Esta operación cuenta con una forma distinta de trabajo a las definidas anteriormente. En este caso el cliente solicita el listado de archivos al servidor. Este último devuelve un único string que contiene el nombre de los archivos, o error en caso de existir.

En este caso la interacción es única, por lo que esta operación se define utilizando el modelo **Simple RPC**.

Para todo esto definimos dos tipos de mensajes: `ListRequest` y `ListResponse` con los datos necesarios, y el servicio `List` de la siguiente forma:

```proto
syntax = "proto3";

package ftp;

message ListRequest {
    string name = 1;
    bool list = 2;
}

message ListResponse {
    string paths = 3;
}

service Operations {
    rpc List(ListRequest) returns (ListResponse);
}
```

### Definición del Servicio

El primer paso es el de definir el servicio y tipos de los mensajes de _request_ y _response_ gRPC utilizando _protocol buffers_.
Para esto definimos el archivo `ftp.proto` con el siguiente contenido:

```proto
syntax = "proto3";

package ftp;

message ReadRequest {
    string name = 1;
    int64 pos = 2;
    int64 bytes = 3;
}

message WriteRequest {
    string name = 1;
    string data = 2;
    int64 checksum = 3;
}

message ListRequest {
    string name = 1;
    bool list = 2;
}

message ReadResponse {
    string name = 1;
    string data = 2;
    bool continue_reading = 3;
}

message WriteResponse {}

message ListResponse {
    string paths = 3;
}

service Operations {
    rpc Read(ReadRequest) returns (stream ReadResponse);
    rpc Write(stream WriteRequest) returns (WriteResponse);
    rpc List(ListRequest) returns (ListResponse);
}
```

### Generar código de Cliente y Servidor

A continuación, debemos generar las interfaces de cliente y servidor de gRPC a partir de nuestra definición de servicio `.proto`. Hacemos esto usando el compilador `protoc` con un complemento especial de GO para gRPC ejecutando el siguiente comando:

```sh
$ protoc --go_out=plugins=grpc:. ./ftp.proto
```

Esto genera el archivo `ftp.pb.go` el cual contiene todo el código del protocol buffer para completar, serializar y recuperar tipos de mensajes de solicitud y respuesta, así como también los tipos necesarios para crear las implementaciones del cliente y del servidor.

### Crear el Servidor

Para esto se deben resolver dos aspectos para hacer que nuestro servicio `FTP` haga su trabajo:

- Implementar la interfaz de servidor generada a partir de nuestra definición de servicio: hacer la implementación real de nuestro servicio.
- Ejecutar un servidor gRPC para escuchar las solicitudes de los clientes y enviarlas a la implementación del servicio adecuada.

Para esto creamos un archivo `.go` con el path `./src/gRPC/ftp/server/main.go` con el contenido que se explica a continuación

#### Implementar Servicio FTP

Para esto creamos el tipo `server` que cuente con un método por cada servicio gRPC que deseamos implementar:

```go
package main

import (
	"context"

    // módulo definido en el archivo ftp.pb.go
	ftp ".."

    . . .
)

type server struct{
    . . .
}

func (s *server) Read(req *ftp.ReadRequest, stream ftp.Operations_ReadServer) error {
	. . .
}

func (s *server) Write(stream ftp.Operations_WriteServer) error {
    . . .
}

func (s *server) List(ctx context.Context, req *ftp.ListRequest) (res *ftp.ListResponse, err error) {
    . . .
}
```

#### Empezando el servidor

Para esto agregamos una función `main` para inicializar el servidor y hacer que este defina el servicio a ejecutar a partir de una _request_.

```go
package main

import (
	"context"

    // módulo definido en el archivo ftp.pb.go
	ftp ".."

	"net"

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
	. . .
}

func (s *server) Write(stream ftp.Operations_WriteServer) error {
    . . .
}

func (s *server) List(ctx context.Context, req *ftp.ListRequest) (res *ftp.ListResponse, err error) {
    . . .
}
```

### Crear el Cliente

Para esto simplemente debemos llamar a los métodos de los servicios generados por el compilador `protoc` a partir del archivo `ftp.proto`.

Teniendo en mente eso es que nuestro cliente FTP en GO deberá crear una conexión con el servidor de la siguiente forma:

```go
package main

import (
	"context"

    // módulo definido en el archivo ftp.pb.go
	ftp ".."

	"google.golang.org/grpc"
)

func main() {
    conn, err := grpc.Dial("localhost:4444", grpc.WithInsecure())
    if err != nil {
        panic(err)
    }

    client := ftp.NewOperationsClient(conn)
}
```

y ejecutar los métodos dependiendo del modo de ejecución.

#### Read (Server Side Stream)

```go
// send request
req := ftp.ReadRequest{
    Name:  ". . .",
    Pos:   0,
    Bytes: 4096,
}

// Create a stream channel
stream, err := client.Read(context.Background(), &req)
done := make(chan bool)
if err != nil {
    log.Fatalf("open stream error %v", err)
}

for {
    // receive data
    res, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatalf("can not receive %v", err)
    }

    // do something with data
    . . .
}

// close chan
close(done)
```

#### Write (Client Side Streaming)

```go
// Create a stream channel
stream, err := client.Write(context.Background())
done := make(chan bool)
if err != nil {
    log.Fatalf("open stream error %v", err)
}

for {
    // save buffer data
    dataVal := ". . ."

    if (someCondition) {
        break;
    }

    // send request
    req := ftp.WriteRequest{
        Name:     ". . .",
        Data:     dataVal,
        Checksum: checksum(buf),
    }
    if err := stream.Send(&req); err != nil {
        log.Fatalf("can not send %v", err)
    }
}

// closing receive stream
_, err = stream.CloseAndRecv()
if err == io.EOF {
    close(done)
}
if err != nil {
    log.Fatalf("can not receive %v", err)
}
```

#### List (Simple RPC)

```go
req := &ftp.ListRequest{Name: ". . .", List: true}

res, err := client.List(context.Background(), req)
if err != nil {
    log.Fatalf("Error when calling Add: %s", err)
}
log.Printf("%s", res.Paths)
```

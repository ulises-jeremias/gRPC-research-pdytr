Existen dos formas de ejecutar cada una de las aplicaciones cliente/servidor gRPC que se encuentran en este [repositorio](https://github.com/ulises-jeremias/gRPC-research-pdytr).

## Utilizando Docker (Recomendado)

Para esto creamos un setup basado en Docker para la ejecución de las distintas aplicaciones.

### Dependencias

- [Docker](https://www.docker.com/)
- [Bash](https://www.gnu.org/software/bash/)

### Ejecución del Ambiente

Con el objetivo de simplificar su uso creamos un script en bash que permite abstrar la creación y ejecución de los ambientes docker. Para esto de debe ejecuctar el siguiente comando en la raiz del repositorio:

```sh
$ ./bin/grpc --app=<app> --root
```

donde,

```
<app> = simp | ftp | user_lookup | request_stream | response_stream | request_and_response_stream
```

_NOTA: Ejecutar `./bin/grpc -h` para conocer más sobre los flags._

Esto ejecutará un contenedor docker en modo interactivo con un volumen en el directorio de la aplicación elegida.

<p align="center">
  <img src="https://raw.githubusercontent.com/ulises-jeremias/gRPC-research-pdytr/master/static/terminal.png">
</p>

### Ejecución de las Aplicaciones

Una vez que el contenedor se esté ejecutando, podrá ejecutar las aplicaciones. Todas tienen exactamente el mismo esquema de ejecución:

1. Compilar el archivo `./<app>.proto`
2. Ejecutar el servidor en background
3. Ejecutar el cliente (*En algunos casos, como el del cliente FTP, se pueden necesitar flags específicos*)

A continuación se muestran los comandos a utilizar para completar cada paso:

_NOTA: Modificar `<app>` por el ejemplo que se está ejecutando._

```sh
# generate proto file
~> protoc --go_out=plugins=grpc:. ./<app>.proto

# execute server in background
~> go run ./server/main.go &

# execute client
~> go run ./client/main.go
```

* * *

## Ejecución en Ambiente Local

Esta opción requiere algunos pasos adicionales respecto a la opción anterior. Para esto se deberán instalar cada una de las dependencias de gRPC y ejecutar los ejemplos localmente. Recomendamos que ante cualquier problema en el setup del ambiente, se opte por utilizar la [opción anterior con docker](#utilizando-docker-recomendado).

### Dependencias

A continuación se muestran las dependencias necesarias para ejecutar las aplicaciones.

- [GO](https://golang.org/), cualquiera de las tres últimas [versiones principales de GO](https://golang.org/doc/devel/release.html).
  
  Para obtener instrucciones de instalación, consulte la [guía de introducción](https://golang.org/doc/install) de Go.

  _NOTA: Esta guía funciona con la versión 1.16.5 que es la última hasta la fecha_

- Compilador de búfer de protocolo, `protoc`, versión 3.

  Para obtener instrucciones de instalación, consulte [Instalación del compilador](https://grpc.io/docs/protoc-installation/) de protocol buffer.

  _NOTA: Esta guía funciona con la versión 3.15.8 que es la última hasta la fecha_

- GO Plugins para el compilador `protoc`:

  Instale los complementos del compilador `protoc` para GO con los siguientes comandos:

  ```sh
  $ go env -w GO111MODULE=auto # enable module mode
  $ go get -u google.golang.org/grpc
  $ go get -u github.com/golang/protobuf/protoc-gen-go
  ```

- Actualizar la variable `PATH` para que el compilador `protoc` pueda encontrar los plugins:

  ```sh
  $ export PATH="$PATH:$(go env GOPATH)/bin"
  ```

### Ejecución de las Aplicaciones

Una vez instaladas las dependencias, podemos proceder a ejecutar las aplicaciones. Para esto se debe seguir el siguiente ejemplo de comandos. Todas tienen exactamente el mismo esquema de ejecución:

1. Compilar el archivo `./<app>.proto`
2. Ejecutar el servidor en background
3. Ejecutar el cliente (*En algunos casos, como el del cliente FTP, se pueden necesitar flags específicos*)

A continuación se muestran los comandos a utilizar para completar cada paso:

_NOTA: Modificar `<app>` por el ejemplo que se está ejecutando._

```sh
# move to app dir
$ cd ./src/gRPC/<app>

# generate proto file
$ protoc --go_out=plugins=grpc:. ./<app>.proto

# execute server in background
$ go run ./server/main.go &

# execute client
$ go run ./client/main.go
```

donde,

```
<app> = simp | ftp | user_lookup | request_stream | response_stream | request_and_response_stream
```

# Ejecución de las aplicaciones

Existen dos formas de ejecutar cada una de las aplicaciones cliente/servidor gRPC que se encuentran en este [repositorio](https://github.com/ulises-jeremias/gRPC-research-pdytr).

- [Ejecución de las aplicaciones](#ejecución-de-las-aplicaciones)
  - [Utilizando Docker (Recomendado)](#utilizando-docker-recomendado)
    - [Dependencias](#dependencias)
    - [Ejecución](#ejecución)
      - [Ejecución de los Ejemplos](#ejecución-de-los-ejemplos)
  - [Sin Docker](#sin-docker)

## Utilizando Docker (Recomendado)

Para esto creamos un setup basado en Docker para la ejecución de las distintas aplicaciones.

### Dependencias

- [Docker](https://www.docker.com/)
- [Bash](https://www.gnu.org/software/bash/)

### Ejecución

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
  <img src="./static/terminal.png">
</p>

#### Ejecución de los Ejemplos

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

## Sin Docker

TODO

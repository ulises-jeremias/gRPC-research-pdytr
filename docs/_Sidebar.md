
- [Ejecución de las aplicaciones](./GRPC_APP_RUN.md)
  - [Utilizando Docker (Recomendado)](#utilizando-docker-recomendado)
    - [Dependencias](#dependencias)
    - [Ejecución del Ambiente](#ejecución-del-ambiente)
    - [Ejecución de las Aplicaciones](#ejecución-de-las-aplicaciones-1)
  - [Ejecución en Ambiente Local](#ejecución-en-ambiente-local)
    - [Dependencias](#dependencias-1)
    - [Ejecución de las Aplicaciones](#ejecución-de-las-aplicaciones-2)
- [Implementación de Cliente/Servidor FTP con gRPC](./GRPC_APP_IMPL.md)
  - [Diseño](#diseño)
    - [Read](#read)
    - [Write](#write)
    - [List](#list)
  - [Definición del Servicio](#definición-del-servicio)
  - [Generar código de Cliente y Servidor](#generar-código-de-cliente-y-servidor)
  - [Crear el Servidor](#crear-el-servidor)
    - [Implementar Servicio FTP](#implementar-servicio-ftp)
    - [Empezando el servidor](#empezando-el-servidor)
  - [Crear el Cliente](#crear-el-cliente)
    - [Read (Server Side Stream)](#read-server-side-stream)
    - [Write (Client Side Streaming)](#write-client-side-streaming)
    - [List (Simple RPC)](#list-simple-rpc)
# Go Key-Value Store

This project implements a simple  Key-Value store system consisting of two Go services that communicate via gRPC:

1. **REST API Server (Frontend)**: Exposes a JSON REST API for public interaction.
2. **gRPC Key-Value Store Server (Backend)**: Implements the core Key-Value storage logic.

Both services are designed to run in separate Docker containers.

## Table of Contents

- [Project Structure](#project-structure)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Building the Services](#building-the-services)
- [Running the Services](#running-the-services)
- [Testing the Services](#testing-the-services)
- [Extensibility and Modularity](#extensibility-and-modularity)

## Project Structure

```
.
├── cmd/
│   ├── kvstore-grpc-server/  # Go source for the gRPC server
│   │   └── main.go
│   └── kvstore-rest-server/  # Go source for the REST API server
│       └── main.go
├── docker/
│   ├── grpc-server/          # Dockerfile for gRPC server
│   │   └── Dockerfile
│   └── rest-server/          # Dockerfile for REST API server
│       └── Dockerfile
├── proto/                    # Protocol Buffers definitions
│   └── kvstore.proto
└── README.md
```

## Features

The Key-Value store supports the following operations through its JSON REST API:

- **Store**: Store a value at a given key.
- **Retrieve**: Retrieve the value for a given key.
- **Delete**: Delete a given key.

## Prerequisites

Before you begin, ensure you have the following installed:

- [Go](https://golang.org/doc/install) (version 1.24 or higher)
- [Docker](https://docs.docker.com/get-docker/)
- [Git](https://git-scm.com/downloads) (to clone the repository)

## Building the Services

First, clone this repository:

```sh
git clone https://github.com/mbasha/grpc-kv-store.git
cd grpc-kv-store
```

Next, generate the gRPC code from the `.proto` file. You'll need `protoc` and `protoc-gen-go`, `protoc-gen-go-grpc`.

**Install protobuf compiler if you haven't already:**

- On macOS:  
  ```sh
  brew install protobuf
  ```
- On Linux:  
  ```sh
  sudo apt-get install protobuf-compiler
  ```

**Install Go plugins for gRPC:**
```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

**Ensure your GOPATH/bin is in your PATH:**
```sh
export PATH=$PATH:$(go env GOPATH)/bin
```

**Generate Go code from proto:**
```sh
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/kvstore.proto
```

Now, build the Docker images for both services:

**Build gRPC server image:**
```sh
docker build -t kvstore-grpc-server -f docker/grpc-server/Dockerfile .
```

**Build REST API server image:**
```sh
docker build -t kvstore-rest-server -f docker/rest-server/Dockerfile .
```

## Running the Services

You need to run the gRPC server first, and then the REST API server, linking them via Docker's networking.
1. **Create the custom network bridge:**
    ```sh
      docker network create kvstore-net
    ```

2. **Run the gRPC Server:**
    ```sh
    docker run -d --name kvstore-grpc-server --network kvstore-net kvstore-grpc-server
    ```

3. **Run the REST API Server:**
    ```sh
    docker run -d --name kvstore-rest-server -p 8080:8080 --network kvstore-net kvstore-rest-server
    ```
    - `-p 8080:8080`: Maps port 8080 on your host to port 8080 in the container, making the REST API accessible from your machine.

You can verify that both containers are running:

```sh
docker ps
```

You should see both `kvstore-grpc-server` and `kvstore-rest-server` listed.

## Testing the Services

Once both services are running, you can test them using `curl`. The REST API server is accessible on `http://localhost:8080`.

### 1. Store a Value (POST)

```sh
curl -X POST -H "Content-Type: application/json" -d '{"value": "world"}' http://localhost:8080/kv/hello
```

Expected Output:
```json
{"success":true,"message":"Key 'hello' stored successfully"}
```

### 2. Retrieve a Value (GET)

```sh
curl http://localhost:8080/kv/hello
```

Expected Output:
```json
{"success":true,"key":"hello","value":"world"}
```

Try retrieving a non-existent key:

```sh
curl http://localhost:8080/kv/nonexistent
```

Expected Output:
```json
{"success":false,"message":"Key 'nonexistent' not found"}
```

### 3. Delete a Key (DELETE)

```sh
curl -X DELETE http://localhost:8080/kv/hello
```

Expected Output:
```json
{"success":true,"message":"Key 'hello' deleted successfully"}
```

Try deleting a non-existent key:

```sh
curl -X DELETE http://localhost:8080/kv/nonexistent
```

Expected Output:
```json
{"success":false,"message":"Key 'nonexistent' not found"}
```

### Clean Up

To stop and remove the Docker containers:

```sh
docker stop kvstore-rest-server kvstore-grpc-server
docker rm kvstore-rest-server kvstore-grpc-server
```
# Movie API

## Getting started

### Protocol Buffers

Follow the Protocol Buffers compiler [install instructions](https://grpc.io/docs/protoc-installation/). Install the Protocol Buffers compiler plugins:

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Update your `PATH` so that the `protoc` (Protocol Buffers compiler) can find the plugins:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

Generate the protobuf code: 

```sh
protoc -I=api --go_out=. movie.proto
```

### Docker

Start a `hashicorp/consul` container:

```sh
docker run -d -p 8500:8500 -p 8600:8600/udp --name=dev-consul hashicorp/consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0
```
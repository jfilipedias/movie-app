# Movie API

## Getting started

Install the Protocol Buffers compiler:

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

Generate the protobuf code: 

```sh
protoc -I=api --go_out=. movie.proto
```

Start a `hashicorp/consul` container:

```sh
docker run -d -p 8500:8500 -p 8600:8600/udp --name=dev-consul hashicorp/consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0
```
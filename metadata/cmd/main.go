package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jfilipedias/movie-app/grpc/gen"
	"github.com/jfilipedias/movie-app/metadata/internal/handler"
	"github.com/jfilipedias/movie-app/metadata/internal/repository/memory"
	"github.com/jfilipedias/movie-app/metadata/internal/service"
	"github.com/jfilipedias/movie-app/pkg/discovery"
	"github.com/jfilipedias/movie-app/pkg/discovery/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var serviceName = "metadata"

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "API handler port")
	flag.Parse()
	log.Printf("Starting the movie metadata service on port %d\n", port)

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	hostPort := fmt.Sprintf("localhost:%d", port)
	if err = registry.Register(ctx, instanceID, serviceName, hostPort); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(serviceName, instanceID); err != nil {
				log.Printf("Failed to report healthy state:  %v\n", err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	defer registry.Deregister(ctx, serviceName, instanceID)

	repo := memory.NewRepository()
	svc := service.NewMetadataService(repo)
	h := handler.NewGrpcHandler(svc)

	lis, err := net.Listen("tcp", hostPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMetadataServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}

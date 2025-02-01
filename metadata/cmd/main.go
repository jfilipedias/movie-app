package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	httpHandler "github.com/jfilipedias/movie-app/metadata/internal/handler/http"
	"github.com/jfilipedias/movie-app/metadata/internal/repository/memory"
	"github.com/jfilipedias/movie-app/metadata/internal/service/metadata"
	"github.com/jfilipedias/movie-app/pkg/discovery"
	"github.com/jfilipedias/movie-app/pkg/discovery/consul"
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
	svc := metadata.NewService(repo)
	h := httpHandler.New(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /metadata", h.GetMetadataByID)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		panic(err)
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jfilipedias/movie-app/pkg/discovery"
	"github.com/jfilipedias/movie-app/pkg/discovery/consul"
	httphandler "github.com/jfilipedias/movie-app/rating/internal/handler/http"
	"github.com/jfilipedias/movie-app/rating/internal/repository/memory"
	"github.com/jfilipedias/movie-app/rating/internal/service/rating"
)

var serviceName = "rating"

func main() {
	var port int
	flag.IntVar(&port, "port", 8082, "API handler port")
	flag.Parse()
	log.Printf("Starting the movie rating service on port %d\n", port)

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
	svc := rating.NewService(repo)
	h := httphandler.New(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /rating", h.GetAggregatedRating)
	mux.HandleFunc("PUT /rating", h.PutRating)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		panic(err)
	}
}

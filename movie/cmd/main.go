package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	metadatagateway "github.com/jfilipedias/movie-app/movie/internal/gateway/metadata/http"
	ratinggateway "github.com/jfilipedias/movie-app/movie/internal/gateway/rating/http"
	"github.com/jfilipedias/movie-app/movie/internal/handler"
	"github.com/jfilipedias/movie-app/movie/internal/service"
	"github.com/jfilipedias/movie-app/pkg/discovery"
	"github.com/jfilipedias/movie-app/pkg/discovery/consul"
)

var serviceName = "movie"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()
	log.Printf("Starting the movie service on port %d\n", port)

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

	metadataGateway := metadatagateway.NewGateway(registry)
	ratingGateway := ratinggateway.NewGateway(registry)
	svc := service.NewMovieService(metadataGateway, ratingGateway)
	h := handler.NewHttpHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /movie", h.GetMovieDetails)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		panic(err)
	}
}

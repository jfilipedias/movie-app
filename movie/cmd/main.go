package main

import (
	"log"
	"net/http"

	"github.com/jfilipedias/movie-app/movie/internal/controller/movie"
	metadatagateway "github.com/jfilipedias/movie-app/movie/internal/gateway/metadata/http"
	ratinggateway "github.com/jfilipedias/movie-app/movie/internal/gateway/rating/http"
	httphandler "github.com/jfilipedias/movie-app/movie/internal/handler/http"
)

func main() {
	log.Print("Starting the movie service")
	metadataGateway := metadatagateway.NewGateway("localhost:8081")
	ratingGateway := ratinggateway.NewGateway("localhost:8082")
	ctrl := movie.New(metadataGateway, ratingGateway)
	h := httphandler.New(ctrl)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /movie", h.GetMovieDetails)
	if err := http.ListenAndServe(":8083", mux); err != nil {
		panic(err)
	}
}

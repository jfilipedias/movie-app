package main

import (
	"log"
	"net/http"

	"github.com/jfilipedias/movie-app/rating/internal/controller/rating"
	httphandler "github.com/jfilipedias/movie-app/rating/internal/handler/http"
	"github.com/jfilipedias/movie-app/rating/internal/repository/memory"
)

func main() {
	log.Println("Starting the rating service")
	repo := memory.New()
	ctrl := rating.New(repo)
	h := httphandler.New(ctrl)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /rating", h.GetAggregatedRating)
	mux.HandleFunc("PUT /rating", h.PutRating)
	if err := http.ListenAndServe(":8082", mux); err != nil {
		panic(err)
	}
}

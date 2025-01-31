package main

import (
	"log"
	"net/http"

	"github.com/jfilipedias/movie-app/metadata/internal/controller/metadata"
	httpHandler "github.com/jfilipedias/movie-app/metadata/internal/handler/http"
	"github.com/jfilipedias/movie-app/metadata/internal/repository/memory"
)

func main() {
	log.Print("Starting the movie metadata service")
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := httpHandler.New(ctrl)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /metadata", h.GetMetadata)
	if err := http.ListenAndServe(":8081", mux); err != nil {
		panic(err)
	}
}

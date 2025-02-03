package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/jfilipedias/movie-app/movie/internal/service"
)

type HttpHandler struct {
	svc *service.MovieService
}

func NewHttpHandler(svc *service.MovieService) *HttpHandler {
	return &HttpHandler{svc}
}

func (h *HttpHandler) GetMovieDetails(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	details, err := h.svc.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(details); err != nil {
		log.Printf("Response encode error: %v\n", err)
	}
}

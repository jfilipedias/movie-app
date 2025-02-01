package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/jfilipedias/movie-app/movie/internal/service/movie"
)

type Handler struct {
	svc *movie.Service
}

func New(svc *movie.Service) *Handler {
	return &Handler{svc}
}

func (h *Handler) GetMovieDetails(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	details, err := h.svc.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, movie.ErrNotFound) {
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

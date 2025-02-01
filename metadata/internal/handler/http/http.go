package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/jfilipedias/movie-app/metadata/internal/service/metadata"
)

type Handler struct {
	svc *metadata.Service
}

func New(svc *metadata.Service) *Handler {
	return &Handler{svc}
}

func (h *Handler) GetMetadataByID(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m, err := h.svc.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, metadata.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Printf("Controller get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

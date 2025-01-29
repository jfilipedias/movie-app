package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/jfilipedias/movie-app/metadata/internal/controller/metadata"
)

type Handler struct {
	ctrl *metadata.Controller
}

func New(ctrl *metadata.Controller) *Handler {
	return &Handler{ctrl}
}

func (h *Handler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	m, err := h.ctrl.Get(ctx, id)
	if err != nil {
		if errors.Is(err, metadata.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("Response encode error: %v\n", err)
	}
}

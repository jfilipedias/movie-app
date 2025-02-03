package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/jfilipedias/movie-app/metadata/internal/service"
)

type HttpHandler struct {
	svc *service.MetadataService
}

func NewHttpHandler(svc *service.MetadataService) *HttpHandler {
	return &HttpHandler{svc}
}

func (h *HttpHandler) GetMetadataByID(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m, err := h.svc.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
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

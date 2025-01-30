package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/jfilipedias/movie-app/rating/internal/controller/rating"
	model "github.com/jfilipedias/movie-app/rating/pkg"
)

type Handler struct {
	ctrl *rating.Controller
}

func New(ctrl *rating.Controller) *Handler {
	return &Handler{ctrl}
}

func (h *Handler) GetAggregatedRating(w http.ResponseWriter, r *http.Request) {
	recordID := model.RecordID(r.FormValue("id"))
	if recordID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recordType := model.RecordType(r.FormValue("type"))
	if recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	v, err := h.ctrl.GetAggregatedRating(r.Context(), recordID, recordType)
	if err != nil {
		if errors.Is(err, rating.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Printf("Controller get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) PutRating(w http.ResponseWriter, r *http.Request) {
	recordID := model.RecordID(r.FormValue("id"))
	if recordID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recordType := model.RecordType(r.FormValue("type"))
	if recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := model.UserID(r.FormValue("userId"))
	v, err := strconv.ParseFloat(r.FormValue("value"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rating := &model.Rating{UserID: userID, Value: model.RatingValue(v)}
	if err = h.ctrl.PutRating(r.Context(), recordID, recordType, rating); err != nil {
		log.Printf("Controller put error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

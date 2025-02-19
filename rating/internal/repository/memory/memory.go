package memory

import (
	"context"

	"github.com/jfilipedias/movie-app/rating/internal/repository"
	"github.com/jfilipedias/movie-app/rating/pkg/model"
)

type Repository struct {
	data map[model.RecordType]map[model.RecordID][]model.Rating
}

func NewRepository() *Repository {
	return &Repository{data: map[model.RecordType]map[model.RecordID][]model.Rating{}}
}

func (r *Repository) Get(_ context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error) {
	if _, ok := r.data[recordType]; !ok {
		return nil, repository.ErrNotFound
	}

	ratings, ok := r.data[recordType][recordID]
	if !ok || len(ratings) == 0 {
		return nil, repository.ErrNotFound
	}
	return ratings, nil
}

func (r *Repository) Put(_ context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = map[model.RecordID][]model.Rating{}
	}

	r.data[recordType][recordID] = append(r.data[recordType][recordID], *rating)
	return nil
}

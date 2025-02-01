package rating

import (
	"context"
	"errors"

	"github.com/jfilipedias/movie-app/rating/internal/repository"
	model "github.com/jfilipedias/movie-app/rating/pkg"
)

var ErrNotFound = errors.New("ratings not found for a record")

type ratingRepository interface {
	Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error)
	Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

type Service struct {
	repo ratingRepository
}

func NewService(repo ratingRepository) *Service {
	return &Service{repo}
}

func (s *Service) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	ratings, err := s.repo.Get(ctx, recordID, recordType)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return 0, ErrNotFound
		}
		return 0, err
	}

	sum := float64(0)

	for _, r := range ratings {
		sum += float64(r.Value)
	}

	return sum / float64(len(ratings)), nil
}

func (s *Service) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	return s.repo.Put(ctx, recordID, recordType, rating)
}

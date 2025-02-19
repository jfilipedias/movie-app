package service

import (
	"context"
	"errors"

	"github.com/jfilipedias/movie-app/rating/internal/repository"
	"github.com/jfilipedias/movie-app/rating/pkg/model"
)

var ErrNotFound = errors.New("ratings not found for a record")

type ratingRepository interface {
	Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error)
	Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

type ratingIngester interface {
	Ingest(ctx context.Context) (chan model.RatingEvent, error)
}

type RatingService struct {
	repo     ratingRepository
	ingester ratingIngester
}

func NewRatingService(repo ratingRepository, ingester ratingIngester) *RatingService {
	return &RatingService{repo, ingester}
}

func (s *RatingService) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
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

func (s *RatingService) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	return s.repo.Put(ctx, recordID, recordType, rating)
}

func (s *RatingService) StartIngestion(ctx context.Context) error {
	ch, err := s.ingester.Ingest(ctx)
	if err != nil {
		return err
	}

	for e := range ch {
		m := &model.Rating{UserID: e.UserID, Value: e.Value}
		if err = s.PutRating(ctx, model.RecordID(e.RecordID), model.RecordType(e.RecordType), m); err != nil {
			return err
		}
	}

	return nil
}

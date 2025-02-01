package movie

import (
	"context"
	"errors"

	metadatamodel "github.com/jfilipedias/movie-app/metadata/pkg"
	"github.com/jfilipedias/movie-app/movie/internal/gateway"
	model "github.com/jfilipedias/movie-app/movie/pkg"
	ratingmodel "github.com/jfilipedias/movie-app/rating/pkg"
)

var ErrNotFound = errors.New("movie metadata not found")

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadatamodel.Metadata, error)
}

type ratingGateway interface {
	GetAggregatedRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) (float64, error)
	PutRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) error
}

type Service struct {
	metadataGateway metadataGateway
	ratingGateway   ratingGateway
}

func NewService(metadataGateway metadataGateway, ratingGateway ratingGateway) *Service {
	return &Service{metadataGateway, ratingGateway}
}

func (s *Service) Get(ctx context.Context, id string) (*model.MovieDetails, error) {
	metadata, err := s.metadataGateway.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gateway.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	details := &model.MovieDetails{Metadata: *metadata}
	rating, err := s.ratingGateway.GetAggregatedRating(ctx, ratingmodel.RecordID(id), ratingmodel.RecordTypeMovie)
	if err != nil {
		return nil, err
	}

	details.Rating = &rating
	return details, nil
}

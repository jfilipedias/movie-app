package service

import (
	"context"
	"errors"

	metadatamodel "github.com/jfilipedias/movie-app/metadata/pkg/model"
	"github.com/jfilipedias/movie-app/movie/internal/gateway"
	"github.com/jfilipedias/movie-app/movie/pkg/model"
	ratingmodel "github.com/jfilipedias/movie-app/rating/pkg/model"
)

var ErrNotFound = errors.New("movie metadata not found")

type metadataGateway interface {
	GetMetadataByID(ctx context.Context, id string) (*metadatamodel.Metadata, error)
}

type ratingGateway interface {
	GetAggregattedRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) (float64, error)
	PutRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) error
}

type MovieService struct {
	metadataGateway metadataGateway
	ratingGateway   ratingGateway
}

func NewMovieService(metadataGateway metadataGateway, ratingGateway ratingGateway) *MovieService {
	return &MovieService{metadataGateway, ratingGateway}
}

func (s *MovieService) GetMovieDetails(ctx context.Context, id string) (*model.MovieDetails, error) {
	metadata, err := s.metadataGateway.GetMetadataByID(ctx, id)
	if err != nil {
		if errors.Is(err, gateway.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	details := &model.MovieDetails{Metadata: *metadata}
	rating, err := s.ratingGateway.GetAggregattedRating(ctx, ratingmodel.RecordID(id), ratingmodel.RecordTypeMovie)
	if err != nil {
		return nil, err
	}

	details.Rating = &rating
	return details, nil
}

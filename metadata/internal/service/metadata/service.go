package metadata

import (
	"context"
	"errors"

	"github.com/jfilipedias/movie-app/metadata/internal/repository"
	model "github.com/jfilipedias/movie-app/metadata/pkg"
)

var ErrNotFound = errors.New("not found")

type metadataRepository interface {
	Get(ctx context.Context, id string) (*model.Metadata, error)
}

type Service struct {
	repo metadataRepository
}

func NewService(repo metadataRepository) *Service {
	return &Service{repo}
}

func (s *Service) Get(ctx context.Context, id string) (*model.Metadata, error) {
	res, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return res, nil
}

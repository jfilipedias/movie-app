package service

import (
	"context"
	"errors"

	"github.com/jfilipedias/movie-app/metadata/internal/repository"
	"github.com/jfilipedias/movie-app/metadata/pkg/model"
)

var ErrNotFound = errors.New("not found")

type metadataRepository interface {
	Get(ctx context.Context, id string) (*model.Metadata, error)
}

type MetadataService struct {
	repo metadataRepository
}

func NewMetadataService(repo metadataRepository) *MetadataService {
	return &MetadataService{repo}
}

func (s *MetadataService) Get(ctx context.Context, id string) (*model.Metadata, error) {
	res, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return res, nil
}

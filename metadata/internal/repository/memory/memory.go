package memory

import (
	"context"
	"sync"

	"github.com/jfilipedias/movie-app/metadata/internal/model"
	"github.com/jfilipedias/movie-app/metadata/internal/repository"
)

type Memory struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

func New() *Memory {
	return &Memory{data: map[string]*model.Metadata{}}
}

func (r *Memory) Get(_ context.Context, id string) (*model.Metadata, error) {
	r.RLock()
	defer r.RUnlock()

	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return m, nil
}

func (r *Memory) Put(_ context.Context, id string, metadata *model.Metadata) error {
	r.Lock()
	defer r.Unlock()
	r.data[id] = metadata
	return nil
}

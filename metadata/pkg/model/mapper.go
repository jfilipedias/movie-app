package model

import "github.com/jfilipedias/movie-app/grpc/gen"

func MetadataToProto(m *Metadata) *gen.Metadata {
	return &gen.Metadata{
		Id:          m.ID,
		Description: m.Description,
		Title:       m.Title,
		Director:    m.Director,
	}
}

func MetadataFromProto(m *gen.Metadata) *Metadata {
	return &Metadata{
		ID:          m.Id,
		Description: m.Description,
		Title:       m.Title,
		Director:    m.Director,
	}
}

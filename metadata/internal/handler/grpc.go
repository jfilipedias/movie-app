package handler

import (
	"github.com/jfilipedias/movie-app/grpc/gen"
	"github.com/jfilipedias/movie-app/metadata/internal/service"
)

type GrpcRandler struct {
	gen.UnimplementedMetadataServiceServer
	svc *service.MetadataService
}

func NewGrpcHandler(svc *service.MetadataService) *HttpHandler {
	return &HttpHandler{svc: svc}
}

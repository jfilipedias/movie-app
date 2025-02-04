package handler

import (
	"context"

	"github.com/jfilipedias/movie-app/grpc/gen"
	"github.com/jfilipedias/movie-app/movie/internal/service"
)

type GrpcHandler struct {
	gen.UnimplementedMovieServiceServer
	svc *service.MovieService
}

func NewGrpcHandler(svc *service.MovieService) *GrpcHandler {
	return &GrpcHandler{svc: svc}
}

func (h *GrpcHandler) GetMovieDetails(ctx context.Context, r *gen.GetMovieDetailsRequest) (*gen.GetMovieDetailsResponse, error) {
	return nil, nil
}

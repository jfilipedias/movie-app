package handler

import (
	"context"
	"errors"

	"github.com/jfilipedias/movie-app/grpc/gen"
	"github.com/jfilipedias/movie-app/metadata/internal/service"
	"github.com/jfilipedias/movie-app/metadata/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
	gen.UnimplementedMetadataServiceServer
	svc *service.MetadataService
}

func NewGrpcHandler(svc *service.MetadataService) *GrpcHandler {
	return &GrpcHandler{svc: svc}
}

func (h *GrpcHandler) GetMetadataByID(ctx context.Context, r *gen.GetMetadataByIdRequest) (*gen.GetMetadataByIdResponse, error) {
	if r == nil || r.MovieId == "" {
		return nil, status.Error(codes.InvalidArgument, "nil request or empty movie id")
	}

	m, err := h.svc.Get(ctx, r.MovieId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.GetMetadataByIdResponse{Metadata: model.MetadataToProto(m)}, nil
}

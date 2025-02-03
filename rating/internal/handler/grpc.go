package handler

import (
	"context"
	"errors"

	"github.com/jfilipedias/movie-app/grpc/gen"
	"github.com/jfilipedias/movie-app/rating/internal/service"
	"github.com/jfilipedias/movie-app/rating/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
	gen.UnimplementedRatingServiceServer
	svc *service.RatingService
}

func NewGrpcHandler(svc *service.RatingService) *GrpcHandler {
	return &GrpcHandler{svc: svc}
}

func (h *GrpcHandler) GetAggregattedRating(ctx context.Context, r *gen.GetAggregattedRatingRequest) (*gen.GetAggregattedRatingResponse, error) {
	if r == nil || r.RecordId == "" || r.RecordType == "" {
		return nil, status.Error(codes.InvalidArgument, "nil request or nil record id/type")
	}

	v, err := h.svc.GetAggregatedRating(ctx, model.RecordID(r.RecordId), model.RecordType(r.RecordType))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.GetAggregattedRatingResponse{RatingValue: v}, nil
}

func (h *GrpcHandler) PutRating(ctx context.Context, r *gen.PutRattingRequest) (*gen.PutRattingResponse, error) {
	if r == nil || r.RecordId == "" || r.RecordType == "" || r.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "nil request or nil user id or record id/type")
	}

	m := &model.Rating{UserID: model.UserID(r.UserId), Value: model.RatingValue(r.RatingValue)}
	if err := h.svc.PutRating(ctx, model.RecordID(r.RecordId), model.RecordType(r.RecordType), m); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.PutRattingResponse{}, nil
}

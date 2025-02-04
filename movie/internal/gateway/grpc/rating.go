package gateway

import (
	"context"

	"github.com/jfilipedias/movie-app/grpc/gen"
	"github.com/jfilipedias/movie-app/movie/internal/grpcutil"
	"github.com/jfilipedias/movie-app/pkg/discovery"
	"github.com/jfilipedias/movie-app/rating/pkg/model"
)

type RatingGrpcGateway struct {
	registry discovery.Registry
}

func NewRatingGrpcGateway(registry discovery.Registry) *RatingGrpcGateway {
	return &RatingGrpcGateway{registry}
}

func (g *RatingGrpcGateway) GetAggregattedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	client := gen.NewRatingServiceClient(conn)
	in := &gen.GetAggregattedRatingRequest{RecordId: string(recordID), RecordType: string(recordType)}
	res, err := client.GetAggregattedRating(ctx, in)
	if err != nil {
		return 0, err
	}

	return res.RatingValue, nil
}

func (g *RatingGrpcGateway) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := gen.NewRatingServiceClient(conn)
	in := &gen.PutRattingRequest{RecordId: string(recordID), RecordType: string(recordType), RatingValue: int32(rating.Value), UserId: string(rating.UserID)}
	_, err = client.PutRating(ctx, in)
	return err
}

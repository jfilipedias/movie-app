package gateway

import (
	"context"

	"github.com/jfilipedias/movie-app/grpc/gen"
	"github.com/jfilipedias/movie-app/metadata/pkg/model"
	"github.com/jfilipedias/movie-app/movie/internal/grpcutil"
	"github.com/jfilipedias/movie-app/pkg/discovery"
)

type GrpcGateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *GrpcGateway {
	return &GrpcGateway{registry}
}

func (g *GrpcGateway) GetMovieDetails(ctx context.Context, id string) (*model.Metadata, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := gen.NewMetadataServiceClient(conn)
	res, err := client.GetMetadata(ctx, &gen.GetMetadataRequest{MovieId: id})
	if err != nil {
		return nil, err
	}

	return model.MetadataFromProto(res.Metadata), nil
}

package gateway

import (
	"context"

	"github.com/jfilipedias/movie-app/grpc/gen"
	"github.com/jfilipedias/movie-app/metadata/pkg/model"
	"github.com/jfilipedias/movie-app/movie/internal/grpcutil"
	"github.com/jfilipedias/movie-app/pkg/discovery"
)

type MetadataGrpcGateway struct {
	registry discovery.Registry
}

func NewMetadataGrpcGateway(registry discovery.Registry) *MetadataGrpcGateway {
	return &MetadataGrpcGateway{registry}
}

func (g *MetadataGrpcGateway) GetMetadataByID(ctx context.Context, id string) (*model.Metadata, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := gen.NewMetadataServiceClient(conn)
	res, err := client.GetMetadataById(ctx, &gen.GetMetadataByIdRequest{MovieId: id})
	if err != nil {
		return nil, err
	}

	return model.MetadataFromProto(res.Metadata), nil
}

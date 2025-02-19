package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/jfilipedias/movie-app/metadata/pkg/model"
	"github.com/jfilipedias/movie-app/movie/internal/gateway"
	"github.com/jfilipedias/movie-app/pkg/discovery"
)

type MetadataHttpGateway struct {
	registry discovery.Registry
}

func NewMetadataHttpGateway(registry discovery.Registry) *MetadataHttpGateway {
	return &MetadataHttpGateway{registry}
}

func (g *MetadataHttpGateway) GetMovieDetails(ctx context.Context, id string) (*model.Metadata, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s/metadata", addrs[rand.Intn(len(addrs))])
	log.Printf("Calling metadata service. Request: GET %s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", id)
	req.URL.RawQuery = values.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if res.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", res)
	}

	var v *model.Metadata
	if err = json.NewDecoder(res.Body).Decode(v); err != nil {
		return nil, err
	}
	return v, nil
}

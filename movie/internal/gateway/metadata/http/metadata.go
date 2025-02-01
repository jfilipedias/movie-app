package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	model "github.com/jfilipedias/movie-app/metadata/pkg"
	"github.com/jfilipedias/movie-app/movie/internal/gateway"
	"github.com/jfilipedias/movie-app/pkg/discovery"
)

type Gateway struct {
	registry discovery.Registry
}

func NewGateway(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s/rating", addrs[rand.Intn(len(addrs))])
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

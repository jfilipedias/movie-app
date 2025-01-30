package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	model "github.com/jfilipedias/movie-app/metadata/pkg"
	"github.com/jfilipedias/movie-app/movie/internal/gateway"
)

type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{addr}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	req, err := http.NewRequest(http.MethodGet, g.addr+"/metadata", nil)
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

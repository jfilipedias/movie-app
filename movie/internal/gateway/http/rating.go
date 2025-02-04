package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/jfilipedias/movie-app/movie/internal/gateway"
	"github.com/jfilipedias/movie-app/pkg/discovery"
	"github.com/jfilipedias/movie-app/rating/pkg/model"
)

type RatingHttpGateway struct {
	registry discovery.Registry
}

func NewRatingHttpGateway(registry discovery.Registry) *RatingHttpGateway {
	return &RatingHttpGateway{registry}
}

func (g *RatingHttpGateway) GetAggregattedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "rating")
	if err != nil {
		return 0, err
	}

	url := fmt.Sprintf("http://%s/rating", addrs[rand.Intn(len(addrs))])
	log.Printf("Calling rating service. Request: GET %s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", string(recordType))
	req.URL.RawQuery = values.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return 0, gateway.ErrNotFound
	} else if res.StatusCode/100 != 2 {
		return 0, fmt.Errorf("non-2xx response: %v", res)
	}

	var v float64
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (g *RatingHttpGateway) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	addrs, err := g.registry.ServiceAddresses(ctx, "rating")
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/rating", addrs[rand.Intn(len(addrs))])
	log.Printf("Calling rating service. Request: PUT %s", url)

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", fmt.Sprintf("%v", recordType))
	values.Add("userId", string(rating.UserID))
	values.Add("value", fmt.Sprintf("%v", rating.Value))
	req.URL.RawQuery = values.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		return fmt.Errorf("non-2xx response: %v", res)
	}
	return nil
}

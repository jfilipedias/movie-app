package ingester

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jfilipedias/movie-app/rating/pkg/model"
)

type KafkaIngester struct {
	consumer *kafka.Consumer
	topic    string
}

func NewKafkaIngester(addr, groupID, topic string) (*KafkaIngester, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": addr,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	return &KafkaIngester{consumer, topic}, nil
}

func (i *KafkaIngester) Ingest(ctx context.Context) (chan model.RatingEvent, error) {
	return nil, nil
}

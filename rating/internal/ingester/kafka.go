package ingester

import (
	"context"
	"encoding/json"
	"fmt"

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
	if err := i.consumer.SubscribeTopics([]string{i.topic}, nil); err != nil {
		return nil, err
	}

	ch := make(chan model.RatingEvent, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				i.consumer.Close()
			default:
			}

			msg, err := i.consumer.ReadMessage(-1)
			if err != nil {
				fmt.Printf("Consumer error: %v\n", err)
				continue
			}

			var event model.RatingEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				fmt.Printf("Unmarshal error: %v\n", err)
				continue
			}
			ch <- event
		}
	}()

	return ch, nil
}

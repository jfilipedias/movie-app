package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jfilipedias/movie-app/rating/pkg/model"
)

func main() {
	ratingEvents, err := readRatingEvents()
	if err != nil {
		panic(err)
	}

	if err = produceRatingEvents(ratingEvents); err != nil {
		panic(err)
	}
}

func readRatingEvents() ([]model.RatingEvent, error) {
	fileName := "ratingsdata.json"
	fmt.Printf("Reading ratings events from file %s\n", fileName)

	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ratingEvents []model.RatingEvent
	if err := json.NewDecoder(f).Decode(&ratingEvents); err != nil {
		return nil, err
	}

	return ratingEvents, nil
}

func produceRatingEvents(ratingEvents []model.RatingEvent) error {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"boostrap.servers":  "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}

	defer producer.Close()

	for _, event := range ratingEvents {
		encodedEvent, err := json.Marshal(event)
		if err != nil {
			return err
		}

		topic := "ratings"
		message := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(encodedEvent),
		}

		if err = producer.Produce(message, nil); err != nil {
			return err
		}
	}

	timeout := 10 * time.Second
	fmt.Println("Waiting " + timeout.String() + " until all events get produced")
	producer.Flush(int(timeout.Milliseconds()))

	return nil
}

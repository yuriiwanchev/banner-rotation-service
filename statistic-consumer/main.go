package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func main() {
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")

	fmt.Println("Kafka brokers:", kafkaBrokers)
	fmt.Println("Kafka topic:", kafkaTopic)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBrokers},
		Topic:   kafkaTopic,
		GroupID: "consumer-group-id",
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Message received: %s\n", string(msg.Value))
	}
}

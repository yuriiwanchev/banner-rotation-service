package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		Writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *KafkaProducer) PublishMessage(key, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	err := p.Writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("Failed to write messages: %v\n", err)
		return err
	}
	return nil
}

func (p *KafkaProducer) Close() error {
	return p.Writer.Close()
}

package kafka

import (
	"KafkaS3/internal/infrastructure/logger"
	"context"

	"github.com/segmentio/kafka-go"
)

func StartProducer(ctx context.Context, logger logger.Logger) *kafka.Writer {
	topic := "my-topic"

	cfg := kafka.WriterConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	writer := kafka.NewWriter(cfg)

	return writer
}

func StopProducer(writer *kafka.Writer) {
	writer.Close()
}

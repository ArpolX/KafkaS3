package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func StartProducer(ctx context.Context, logger *zap.SugaredLogger) *kafka.Writer {
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

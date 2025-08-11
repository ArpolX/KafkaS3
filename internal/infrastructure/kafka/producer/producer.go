package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	Writer *kafka.Writer
}

const (
	topic = "my-topic"
)

func StartProducer(ctx context.Context, logger *zap.SugaredLogger) *Producer {
	cfg := kafka.WriterConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	writer := kafka.NewWriter(cfg)

	return &Producer{Writer: writer}
}

func StopProducer(writer *kafka.Writer, logger *zap.SugaredLogger) error {
	if err := writer.Close(); err != nil {
		logger.Error("Некорректное завершение producer, метод StopProducer", zap.Error(err))
		return fmt.Errorf("Некорректное завершение producer, метод StopProducer: %w", err)
	}

	return nil
}

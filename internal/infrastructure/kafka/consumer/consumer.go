package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

const (
	groupID = "my-consumer-group"
	topic   = "my-topic"
)

func StartConsumerAndReadMessage(ctx context.Context, logger *zap.SugaredLogger) error {
	cfg := kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   topic,
		GroupID: groupID,
		MaxWait: 1 * time.Second,
	}

	reader := kafka.NewReader(cfg)

	defer func() {
		if err := reader.Close(); err != nil {
			logger.Error("Некорректное завершение consumer, метод StartConsumerAndReadMessage", zap.Error(err))
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			message, err := reader.ReadMessage(ctx)
			if err != nil {
				logger.Error("Ошибка приёма сообщения в consumer", zap.Error(err))
			}
			if message.Value == nil {
				logger.Error("Consumer: получено пустое сообщение")
			}
			logger.Info("Consumer: сообщение получено")
		}
	}
}

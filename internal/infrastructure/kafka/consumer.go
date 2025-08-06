package kafka

import (
	"KafkaS3/internal/infrastructure/logger"
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

const (
	groupID = "my-consumer-group"
)

func StartConsumer(ctx context.Context, logger logger.Logger) {
	topic := "my-topic"

	cfg := kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   topic,
		GroupID: groupID,
	}

	reader := kafka.NewReader(cfg)
	defer reader.Close()

	for {
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			logger.Error("Ошибка приёма сообщения в consumer", zap.Error(err))
		}
		if message.Value == nil {
			logger.Error("Пустое сообщение в consumer")
		}
		logger.Info("Сообщение получено consumer")
	}
}

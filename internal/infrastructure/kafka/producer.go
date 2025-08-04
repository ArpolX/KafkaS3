package kafka

import (
	"KafkaS3/internal/controller"
	"KafkaS3/internal/infrastructure/logger"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func StartProducer(ctx context.Context, logger logger.Logger) {
	topic := "my-topic"

	cfg := kafka.WriterConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	writer := kafka.NewWriter(cfg)
	defer writer.Close()

	j := 0
	for {
		producerData := controller.GenerateFakeData()

		data, err := json.Marshal(producerData)
		if err != nil {
			logger.Error("Ошибка сериализация данных", zap.Error(err))
			return
		}

		err = writer.WriteMessages(ctx,
			kafka.Message{
				Key:   []byte(strconv.Itoa(j)),
				Value: data,
			},
		)
		if err != nil {
			logger.Error("Ошибка приёма сообщения в producer", zap.Error(err))
			return
		}
		time.Sleep(1000 * time.Millisecond)
		j++
	}
}

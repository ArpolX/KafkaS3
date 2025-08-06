package controller

import (
	"KafkaS3/internal/infrastructure/logger"
	"KafkaS3/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/minio/minio-go"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Controller struct {
	KafkaWriter *kafka.Writer
	MinioClient *minio.Client
	Logger      logger.Logger
	Service     service.Service
}

func NewController(writer *kafka.Writer, logger logger.Logger, service service.Service) *Controller {
	return &Controller{
		KafkaWriter: writer,
		Logger:      logger,
		Service:     service,
	}
}

func (c *Controller) DispatchKafka(ctx context.Context) error {
	j := 0
	for {
		producerData := c.Service.GenerateFakeData()

		data, err := json.Marshal(producerData)
		if err != nil {
			c.Logger.Error("Ошибка сериализация данных", zap.Error(err))
			return fmt.Errorf("Ошибка сериализация данных: %w", err)
		}

		err = c.KafkaWriter.WriteMessages(ctx,
			kafka.Message{
				Key:   []byte(strconv.Itoa(j)),
				Value: data,
			},
		)
		if err != nil {
			c.Logger.Error("Ошибка отправки сообщения producer", zap.Error(err))
			return fmt.Errorf("Ошибка отправки сообщения producer: %w", err)
		}
		time.Sleep(1000 * time.Millisecond)
		j++
	}
}

func (c *Controller) UploadImage(ctx context.Context) error {

}

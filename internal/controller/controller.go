package controller

import (
	"KafkaS3/internal/infrastructure/s3"
	"KafkaS3/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Controller struct {
	KafkaWriter *kafka.Writer
	S3          *s3.S3
	Logger      *zap.SugaredLogger
	Service     service.Service
}

func NewController(writer *kafka.Writer, s3 *s3.S3, logger *zap.SugaredLogger, service service.Service) *Controller {
	return &Controller{
		KafkaWriter: writer,
		Logger:      logger,
		Service:     service,
		S3:          s3,
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

func (c *Controller) UploadImage(ctx context.Context, filePath, objectName string) (minio.UploadInfo, error) {
	info, err := c.S3.Minio.FPutObject(ctx, c.S3.Bucket, objectName, filePath, minio.PutObjectOptions{
		ContentType: "image/png",
	})
	if err != nil {
		c.Logger.Error("Ошибка отправки изображения в s3, метод UploadImage", zap.Error(err))
		return minio.UploadInfo{}, fmt.Errorf("Ошибка отправки изображения в s3, метод UploadImage: %w", err)
	}
	c.Logger.Info("Изображение отправлено")

	return info, nil
}

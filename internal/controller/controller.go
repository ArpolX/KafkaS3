package controller

import (
	p "KafkaS3/internal/infrastructure/kafka/producer"
	"KafkaS3/internal/infrastructure/s3"
	"KafkaS3/internal/service"

	"go.uber.org/zap"
)

type Controller struct {
	Producer *p.Producer
	S3       *s3.S3
	Logger   *zap.SugaredLogger
	Service  service.Service
}

func NewController(producer *p.Producer, s3 *s3.S3, logger *zap.SugaredLogger, service service.Service) *Controller {
	return &Controller{
		Producer: producer,
		Logger:   logger,
		Service:  service,
		S3:       s3,
	}
}

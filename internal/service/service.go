package service

import (
	"KafkaS3/internal/dto"

	"go.uber.org/zap"
)

type ServiceImpl struct {
	Logger *zap.SugaredLogger
}

type Service interface {
	GenerateFakeData() []*dto.ProducerData
}

func NewServiceImpl(logger *zap.SugaredLogger) Service {
	return &ServiceImpl{
		Logger: logger,
	}
}

package service

import (
	"KafkaS3/internal/entity"

	"go.uber.org/zap"
)

type ServiceImpl struct {
	Logger *zap.SugaredLogger
}

type Service interface {
	GenerateFakeData() []*entity.FakeDataUser
}

func NewServiceImpl(logger *zap.SugaredLogger) Service {
	return &ServiceImpl{
		Logger: logger,
	}
}

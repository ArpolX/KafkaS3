package service

import (
	"KafkaS3/internal/dto"
	"KafkaS3/internal/infrastructure/logger"

	"github.com/brianvoe/gofakeit"
)

type ServiceImpl struct {
	Logger logger.Logger
}

type Service interface {
	GenerateFakeData() []*dto.ProducerData
}

func NewServiceImpl(logger logger.Logger) Service {
	return &ServiceImpl{
		Logger: logger,
	}
}

func (s *ServiceImpl) GenerateFakeData() []*dto.ProducerData {
	producerData := []*dto.ProducerData{}
	for i := 0; i < 100; i++ {
		data := dto.ProducerData{
			Id:        i,
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			City:      gofakeit.City(),
			Phone:     gofakeit.Phone(),
		}

		producerData = append(producerData, &data)
	}

	return producerData
}

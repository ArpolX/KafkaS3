package service

import (
	"KafkaS3/internal/dto"

	"github.com/brianvoe/gofakeit"
)

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

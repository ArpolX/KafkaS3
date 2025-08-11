package service

import (
	"KafkaS3/internal/entity"

	"github.com/brianvoe/gofakeit"
)

func (s *ServiceImpl) GenerateFakeData() []*entity.FakeDataUser {
	producerData := []*entity.FakeDataUser{}
	for i := 0; i < 100; i++ {
		data := entity.FakeDataUser{
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

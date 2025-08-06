package s3

import (
	"KafkaS3/internal/infrastructure/logger"
	"fmt"

	"github.com/minio/minio-go"
	"go.uber.org/zap"
)

const (
	endpoint        = "minio:9090"
	accessKeyId     = "admin"
	secretAccessKey = "password"
	useTLS          = false
)

func NewS3Client(logger logger.Logger) (*minio.Client, error) {
	minioClient, err := minio.New(endpoint, accessKeyId, secretAccessKey, useTLS)
	if err != nil {
		logger.Fatal("Ошибка в создании s3 хранилища, метод NewS3Client", zap.Error(err))
		return nil, fmt.Errorf("Ошибка в создании s3 хранилища, метод NewS3Client: %w", err)
	}

	return minioClient, nil
}

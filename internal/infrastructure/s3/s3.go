package s3

import (
	"KafkaS3/internal/config"
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

const (
	useTLS = false
)

type S3 struct {
	Minio  *minio.Client
	Bucket string
}

func NewS3Client(ctx context.Context, l *zap.SugaredLogger, cfg *config.Config) (*S3, error) {
	log.Println(cfg.S3.Endpoint)
	client, err := minio.New(cfg.S3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3.AccessKeyId, cfg.S3.SecretAccessKey, ""),
		Secure: useTLS,
	})
	if err != nil {
		l.Fatal("Ошибка в создании s3 хранилища, метод NewS3Client", zap.Error(err))
		return nil, fmt.Errorf("Ошибка в создании s3 хранилища, метод NewS3Client: %w", err)
	}

	if err := ensureBucket(ctx, client, cfg.S3.Bucket); err != nil {
		l.Error("Ошибка в проверке бакетов s3, метод NewS3Client", zap.Error(err))
		return nil, fmt.Errorf("Ошибка в проверке бакетов s3, метод NewS3Client: %w", err)
	}

	return &S3{Minio: client, Bucket: cfg.S3.Bucket}, nil
}

func ensureBucket(ctx context.Context, client *minio.Client, bucket string) error {
	ok, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("Ошибка в проверке наличия бакета в s3, метод ensureBucket: %w", err)
	}

	if !ok {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("Ошибка в создании бакета в s3, метод ensureBucket: %w", err)
		}
	}

	return nil
}

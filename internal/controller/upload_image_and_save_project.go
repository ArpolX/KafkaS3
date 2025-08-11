package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

const (
	savePath     = "/app/download_s3_image" // куда сохраняем из s3
	imageDirPath = "/app/image"             // откуда берём
)

func (c *Controller) UploadImageAndSaveProject(ctx context.Context) error {
	files, err := os.ReadDir(imageDirPath)
	if err != nil {
		c.Logger.Error("Ошибка чтения изображений", zap.Error(err))
		return fmt.Errorf("Ошибка чтения изображений: %w", err)
	}

	for _, file := range files {
		filePath := fmt.Sprintf("%v/%v", imageDirPath, file.Name())

		contentType, err := detectContentType(filePath)
		if err != nil {
			c.Logger.Error("Ошибка в определении типа файла", zap.Error(err))
			return fmt.Errorf("Ошибка в определении типа файла: %w", err)
		}

		// сохранили
		info, err := c.S3.Minio.FPutObject(ctx, c.S3.Bucket, file.Name(), filePath, minio.PutObjectOptions{
			ContentType: contentType,
		})
		if err != nil {
			c.Logger.Error("Ошибка отправки изображения в s3, метод UploadImage", zap.Error(err))
			return fmt.Errorf("Ошибка отправки изображения в s3, метод UploadImage: %w", err)
		}
		c.Logger.Info("Изображение отправлено, информация:", info)

		// взяли
		if err := getImageAndSaveInProject(ctx, c.S3.Minio, c.Logger, c.S3.Bucket, file.Name()); err != nil {
			return fmt.Errorf("Ошибка в загрузке изображения: %w", err)
		}
		c.Logger.Info("Изображение загружено в папку download_s3_image")
	}

	return nil
}

func getImageAndSaveInProject(ctx context.Context, client *minio.Client, l *zap.SugaredLogger, bucket, fileName string) error {
	savePathImage := fmt.Sprintf("%v/%v", savePath, fileName)
	err := client.FGetObject(ctx, bucket, fileName, savePathImage, minio.GetObjectOptions{})
	if err != nil {
		l.Error("Не удалось загрузить изображение, метод GetImageAndSaveInProject", zap.Error(err))
		return fmt.Errorf("Не удалось загрузить изображение, метод GetImageAndSaveInProject: %w", err)
	}

	return nil
}

func detectContentType(filePath string) (string, error) {
	log.Println(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("Ошибка открытия файла, метод detectContentType: %w", err)
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil {
		return "", fmt.Errorf("Ошибка чтения файла, метод detectContentType: %w", err)
	}

	return http.DetectContentType(buf[:n]), nil
}

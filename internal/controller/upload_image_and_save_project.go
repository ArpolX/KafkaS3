package controller

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
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
		fileName := file.Name() + ".gz"

		pr, pw := io.Pipe()

		go func() {
			defer pw.Close()

			file, err := os.Open(filePath)
			if err != nil {
				c.Logger.Error("Ошибка открытия файла, метод UploadImageAndSaveProject", zap.Error(err))
				return
			}
			defer file.Close()

			zip := gzip.NewWriter(pw)
			defer zip.Close()

			_, err = io.Copy(zip, file)
			if err != nil {
				c.Logger.Error("Ошибка при копировании, метод UploadImageAndSaveProject", zap.Error(err))
				return
			}
		}()

		// сохранили
		if err := saveImageInS3(ctx, c.S3.Minio, c.Logger, c.S3.Bucket, fileName, pr); err != nil {
			return fmt.Errorf("Ошибка при сохранении изображения в S3: %w", err)
		}

		// взяли
		if err := saveImageInProject(ctx, c.S3.Minio, c.Logger, c.S3.Bucket, fileName); err != nil {
			return fmt.Errorf("Ошибка при загрузки изображения из S3: %w", err)
		}
		c.Logger.Info("Изображение загружено в папку download_s3_image")
	}

	return nil
}

func saveImageInProject(ctx context.Context, client *minio.Client, l *zap.SugaredLogger, bucket, fileName string) error {
	savePathImage := fmt.Sprintf("%v/%v", savePath, fileName)
	image, err := client.GetObject(ctx, bucket, fileName, minio.GetObjectOptions{})
	if err != nil {
		l.Error("Не удалось загрузить изображение, метод GetImageAndSaveInProject", zap.Error(err))
		return fmt.Errorf("Не удалось загрузить изображение, метод GetImageAndSaveInProject: %w", err)
	}
	defer image.Close()

	zip, err := gzip.NewReader(image)
	if err != nil {
		l.Error("Ошибка декомпрессии", zap.Error(err))
		return fmt.Errorf("Ошибка декомпрессии: %w", err)
	}
	defer zip.Close()

	f, err := os.Create(savePathImage)
	if err != nil {
		l.Error("Ошибка открытия файла сохранения", zap.Error(err))
		return fmt.Errorf("Ошибка открытия файла сохранения: %w", err)
	}
	defer f.Close()

	_, err = io.Copy(f, zip)
	if err != nil {
		l.Error("Ошибка при копировании, метод UploadImageAndSaveProject", zap.Error(err))
		return fmt.Errorf("Ошибка при копировании, метод UploadImageAndSaveProject: %w", err)
	}

	return nil
}

func saveImageInS3(ctx context.Context, client *minio.Client, l *zap.SugaredLogger, bucket, fileName string, file io.Reader) error {
	info, err := client.PutObject(ctx, bucket, fileName, file, -1, minio.PutObjectOptions{
		ContentType:     "image/gzip",
		ContentEncoding: "gzip",
	})
	if err != nil {
		l.Error("Ошибка отправки изображения в s3, метод UploadImage", zap.Error(err))
		return fmt.Errorf("Ошибка отправки изображения в s3, метод UploadImage: %w", err)
	}
	l.Info("Изображение отправлено, информация:", info)
	return nil
}

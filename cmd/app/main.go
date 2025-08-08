package main

import (
	"KafkaS3/internal/app"
	"KafkaS3/internal/config"
	"KafkaS3/internal/infrastructure/logger"
	"context"

	"go.uber.org/zap"
)

func main() {
	l := logger.NewLogger()

	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		l.Fatal("Ошибка чтения конфигурации .env файла", zap.Error(err))
	}

	err = app.Run(ctx, cfg, l)
	if err != nil {
		l.Fatal("Ошибка запуска сервера")
	}
}

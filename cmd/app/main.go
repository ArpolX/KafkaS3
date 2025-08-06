package main

import (
	"KafkaS3/internal/app"
	"KafkaS3/internal/infrastructure/logger"
)

func main() {
	logger := logger.NewLogger()

	err := app.Run(logger)
	if err != nil {
		logger.Fatal("Ошибка запуска сервера")
	}
}

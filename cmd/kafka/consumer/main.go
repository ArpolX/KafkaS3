package main

import (
	"KafkaS3/internal/infrastructure/kafka"
	"KafkaS3/internal/infrastructure/logger"
	"context"
)

func main() {
	ctx := context.Background()

	logger := logger.NewLogger()

	kafka.StartConsumer(ctx, logger)
}

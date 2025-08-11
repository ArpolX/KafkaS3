package main

import (
	c "KafkaS3/internal/infrastructure/kafka/consumer"
	"KafkaS3/internal/infrastructure/logger"
	"context"
)

func main() {
	ctx := context.Background()

	l := logger.NewLogger()

	c.StartConsumerAndReadMessage(ctx, l)
}

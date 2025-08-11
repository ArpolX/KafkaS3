package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func (c *Controller) DispatchKafka(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			producerData := c.Service.GenerateFakeData()

			data, err := json.Marshal(producerData)
			if err != nil {
				c.Logger.Error("Ошибка сериализация данных", zap.Error(err))
				return fmt.Errorf("Ошибка сериализация данных: %w", err)
			}

			err = c.Producer.Writer.WriteMessages(ctx,
				kafka.Message{
					Key:   []byte(producerData[0].FirstName),
					Value: data,
					Time:  time.Now(),
				},
			)
			if err != nil {
				c.Logger.Error("Ошибка отправки сообщения producer", zap.Error(err))
				return fmt.Errorf("Ошибка отправки сообщения producer: %w", err)
			}
			c.Logger.Info("Producer: данные отправлены")

			time.Sleep(7 * time.Second)
		}
	}
}

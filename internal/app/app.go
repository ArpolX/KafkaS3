package app

import (
	"KafkaS3/internal/controller"
	"KafkaS3/internal/infrastructure/kafka"
	"KafkaS3/internal/infrastructure/logger"
	"KafkaS3/internal/service"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func Run(l logger.Logger) error {
	ctx := context.Background()

	writer := kafka.StartProducer(ctx, l)

	Service := service.NewServiceImpl(l)

	Controller := controller.NewController(writer, l, Service)

	go func() {
		err := Controller.DispatchKafka(ctx)
		if err != nil {
			return
		}
	}()

	server := http.Server{}

	list, err := net.Listen("tcp", "app:8080")
	if err != nil {
		l.Fatal("Ошибка открытия соединения", zap.Error(err))
		return fmt.Errorf("Ошибка открытия соединения: %w", err)
	}
	defer list.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Serve(list); err != nil {
			l.Fatal("Ошибка запуска сервера", zap.Error(err))
		}
	}()

	<-stop
	kafka.StopProducer(writer)

	err = server.Shutdown(ctx)
	if err != nil {
		l.Error("Сервер некорректно завершил работу", zap.Error(err))
		return fmt.Errorf("Сервер некорректно завершил работу: %w", err)
	}
	close(stop)
	l.Info("Сервер остановлен Graceful Shutdown")

	return nil
}

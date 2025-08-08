package app

import (
	"KafkaS3/internal/config"
	ctrl "KafkaS3/internal/controller"
	"KafkaS3/internal/infrastructure/kafka"
	"KafkaS3/internal/infrastructure/s3"
	srv "KafkaS3/internal/service"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func Run(ctx context.Context, cfg *config.Config, l *zap.SugaredLogger) error {
	s3, err := s3.NewS3Client(ctx, l, cfg)
	if err != nil {
		l.Errorf("error start s3", zap.Error(err))
		return fmt.Errorf("error start s3: %w", err)
	}

	writer := kafka.StartProducer(ctx, l)

	service := srv.NewServiceImpl(l)

	controller := ctrl.NewController(writer, s3, l, service)

	go func() {
		err := controller.DispatchKafka(ctx)
		if err != nil {
			return
		}
	}()

	go func() {
		info, err := controller.UploadImage(ctx, "/../../image/", "dto.png")
		if err != nil {
			return
		}
		l.Info("s3", info)
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

package app

import (
	"KafkaS3/internal/config"
	ctrl "KafkaS3/internal/controller"
	kafka1 "KafkaS3/internal/infrastructure/kafka/producer"
	"KafkaS3/internal/infrastructure/s3"
	srv "KafkaS3/internal/service"
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func Run(ctx context.Context, cfg *config.Config, l *zap.SugaredLogger) {
	s3, err := s3.NewS3Client(ctx, l, cfg)
	if err != nil {
		l.Errorf("error start s3", zap.Error(err))
	}

	producer := kafka1.StartProducer(ctx, l)

	service := srv.NewServiceImpl(l)

	controller := ctrl.NewController(producer, s3, l, service)

	// начинает слать сообщения producer
	go func() {
		err := controller.DispatchKafka(ctx)
		if err != nil {
			return
		}
	}()

	// отправка и сохранение изображений s3
	go func() {
		err := controller.UploadImageAndSaveProject(ctx)
		if err != nil {
			return
		}
	}()

	server := http.Server{
		Addr:         "app:8888",
		ReadTimeout:  10,
		WriteTimeout: 10,
	}

	list, err := net.Listen("tcp", "app:8888")
	if err != nil {
		l.Fatal("Ошибка открытия соединения", zap.Error(err))
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

	if err := kafka1.StopProducer(producer.Writer, l); err != nil {
		l.Error("error stop kafka1", zap.Error(err))
	}

	err = server.Shutdown(ctx)
	if err != nil {
		l.Error("Сервер некорректно завершил работу", zap.Error(err))
	}
	close(stop)
	l.Info("Сервер остановлен Graceful Shutdown")
}

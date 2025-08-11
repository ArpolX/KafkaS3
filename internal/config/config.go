package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	S3 S3
}

type S3 struct {
	Endpoint        string `env:"ENDPOINT"`
	AccessKeyId     string `env:"ACCESS_KEY_ID"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY"`
	Bucket          string `env:"BUCKET"`
}

func LoadConfig() (*Config, error) {
	projectRoot, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Ошибка получения рабочей директории, метод LoadConfig: %w", err)
	}

	envPath := filepath.Join(projectRoot, ".env")

	var cfg Config
	if err := cleanenv.ReadConfig(envPath, &cfg); err != nil {
		return nil, fmt.Errorf("Ошибка загрузки переменных окружения, метод LoadConfig: %w", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("Ошибка чтения .env файла, метод LoadConfig: %w", err)
	}

	return &cfg, nil
}

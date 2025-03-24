package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSLMode  string
	TimeZone   string
	GRPCPort   string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Предупреждение: не удалось загрузить .env файл, используются системные переменные")
	}

	cfg := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBPort:     os.Getenv("DB_PORT"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
		TimeZone:   os.Getenv("TIMEZONE"),
		GRPCPort:   os.Getenv("GRPC_PORT"),
	}

	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBName == "" || cfg.DBPort == "" {
		return nil, fmt.Errorf("необходимо задать все обязательные параметры для подключения к БД")
	}

	return cfg, nil
}

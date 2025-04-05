package config

import (
	"fmt"
	"os"
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
	// if err := godotenv.Load(); err != nil {
	// 	log.Println("Предупреждение: не удалось загрузить .env файл, используются системные переменные")
	// }
	// this is for local development

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
		return nil, fmt.Errorf("need to set DB_HOST, DB_USER, DB_NAME, DB_PORT")
	}

	return cfg, nil
}

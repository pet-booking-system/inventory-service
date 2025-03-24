package cmd

import (
	"fmt"
	"log"
	"net"

	"invservice/config"
	"invservice/internal/migrations"
	"invservice/internal/repository"
	"invservice/internal/server"
	"invservice/internal/service"

	inventorypb "github.com/azhaxyly/proto-definitions/inventory"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode, cfg.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("Ошибка при миграции: %v", err)
	}

	repo := repository.NewInventoryRepository(db)
	invService := service.NewInventoryService(repo)

	grpcServer := grpc.NewServer()
	inventorypb.RegisterInventoryServiceServer(grpcServer, server.NewInventoryServer(invService))
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Ошибка при открытии порта: %v", err)
	}

	log.Printf("Inventory gRPC сервер запущен на порту %s", cfg.GRPCPort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка запуска gRPC сервера: %v", err)
	}
}

package cmd

import (
	"fmt"
	"log"
	"net"

	"invservice/config"
	"invservice/internal/interceptors"
	"invservice/internal/repository"
	"invservice/internal/server"
	"invservice/internal/service"

	inventorypb "github.com/pet-booking-system/proto-definitions/inventory"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode, cfg.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connect to DB: %v", err)
	}

	repo := repository.NewInventoryRepository(db)
	invService := service.NewInventoryService(repo)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.AuthInterceptor()),
	)
	inventorypb.RegisterInventoryServiceServer(grpcServer, server.NewInventoryServer(invService))
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Error opening port: %v", err)
	}

	log.Printf("Inventory gRPC server running on port: %s", cfg.GRPCPort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error starting gRPC server: %v", err)
	}
}

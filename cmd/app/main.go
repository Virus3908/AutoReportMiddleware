package main

import (
	"context"
	"fmt"
	"log"
	"main/internal/common/interfaces"
	"main/internal/config"
	"main/internal/handlers"
	"main/internal/kafka/consumer"
	"main/internal/kafka/producer"
	"main/internal/logger"
	"main/internal/postgres"
	"main/internal/repositories"
	"main/internal/services"
	"main/internal/storage"
	"net/http"
	"main/internal/middleware"
	
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file found or failed to load it, %s", err.Error())
	}
}

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic("Config file error: " + err.Error())
	}

	serverSettings := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	log, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		panic("failed to init logger: " + err.Error())
	}
	defer log.Sync()
	ctx := context.WithValue(context.Background(), "logger", log)

	db, err := postgres.New(ctx, cfg.DB)
	if err != nil {
		log.Fatal("DB connection error", interfaces.LogField{Key: "Error", Value: err.Error()})
	}
	defer db.Close()

	repo := repositories.New(db.GetPool())

	storage, err := storage.New(ctx, cfg.S3)
	if err != nil {
		log.Fatal("Storage connection error", interfaces.LogField{Key: "Error", Value: err.Error()})
	}

	messageProducer, err := producer.NewProducer(cfg.Producer)
	if err != nil {
		log.Fatal("Producer connection error", interfaces.LogField{Key: "Error", Value: err.Error()})
	}
	defer messageProducer.Close()

	service := services.New(repo, storage, messageProducer, db, true)

	middlewares := []mux.MiddlewareFunc{
		log.LoggingMidleware,
		logger.ContextWithLogger(log),
		middleware.WithCORS,
	}

	router := handlers.New(service, middlewares)
	messageConsumer, err := consumer.NewConsumer(cfg.Consumer, service.Tasks)
	if err != nil {
		log.Fatal("Consumer connection error", interfaces.LogField{Key: "Error", Value: err.Error()})
	}
	defer messageConsumer.Close()

	log.Info("Server is ready",
		interfaces.LogField{Key: "Host", Value: cfg.Server.Host},
		interfaces.LogField{Key: "Port", Value: cfg.Server.Port},
	)

	go messageConsumer.Start(ctx)
	err = http.ListenAndServe(serverSettings, router.GetRouter())
	if err != nil {
		log.Fatal("Server stating error", interfaces.LogField{Key: "Error", Value: err.Error()})
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"main/internal/config"
	"main/internal/handlers"
	"main/internal/kafka/consumer"
	"main/internal/kafka/producer"
	"main/internal/logging"
	"main/internal/postgres"
	"main/internal/repositories"
	"main/internal/services"
	"main/internal/storage"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Config file error: %s", err)
	}

	serverSettings := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)


	db, err := postgres.New(context.Background(), cfg.DB)
	if err != nil {
		log.Fatalf("DB connection error: %s", err)
	}
	defer db.Close()

	repo := repositories.New(db.GetPool())

	storage, err := storage.New(context.Background(), cfg.S3)
	if err != nil {
		log.Fatalf("Storage connection error: %s", err)
	}

	producer, err := producer.NewProducer(cfg.Producer)
	if err != nil {
		log.Fatalf("Producer connection error: %s", err)
	}
	defer producer.Close()

	service := services.New(repo, storage, producer, db, true)

	middlewares := []mux.MiddlewareFunc{
		logging.LoggingMidleware,
	}

	router := handlers.New(service, middlewares)
	consumer, err := consumer.NewConsumer(cfg.Consumer, service.Tasks)
	if err != nil {
		log.Fatalf("Consumer connection error: %s", err)
	}
	defer consumer.Close()

	log.Printf("Server is ready: %s", serverSettings)
	
	go consumer.Start(context.Background()) 
	err = http.ListenAndServe(serverSettings, router.GetRouter())
	if err != nil {
		log.Fatalf("Server stating error: %s", err)
	}
}

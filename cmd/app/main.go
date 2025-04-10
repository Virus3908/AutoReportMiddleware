package main

import (
	"context"
	"fmt"
	"log"
	"main/internal/config"
	"main/internal/handlers"
	"main/internal/kafka"
	"main/internal/logging"
	"main/internal/repositories"
	"main/internal/postgres"
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

	kafkaProducer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		log.Fatalf("Kafka connection error: %s", err)
	}
	defer kafkaProducer.Close()

	service := services.New(repo, storage, kafkaProducer, db, true, cfg.Server.Host, cfg.Server.Port)

	middlewares := []mux.MiddlewareFunc{
		logging.LoggingMidleware,
	}

	router := handlers.New(service, middlewares)

	log.Printf("Server is ready: %s", serverSettings)
	err = http.ListenAndServe(serverSettings, router.GetRouter())
	if err != nil {
		log.Fatalf("Server stating error: %s", err)
	}
}

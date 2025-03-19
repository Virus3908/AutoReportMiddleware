package main

import (
	"context"
	"fmt"
	"log"
	"main/internal/clients"
	"main/internal/config"
	"main/internal/database"
	"main/internal/handlers"
	"main/internal/storage"
	"net/http"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Config file error: %s", err)
	}

	db, err := database.NewDatabase(cfg.DB)
	if err != nil {
		log.Fatalf("DB connection error: %s", err)
	}
	defer db.CloseConnection()

	storage, err := storage.NewStorage(cfg.S3)
	if err != nil {
		log.Fatalf("Storage connection error: %s", err)
	}

	_, err = clients.NewAPIClient(context.Background(), cfg.API)
	if err != nil {
		log.Fatalf("Client connection error: %s", err)
	}

	router := handlers.NewRouter(db, storage)
	router.CreateHandlers()

	serverSettings := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	handlers.SetReady()
	log.Printf("Server is ready: %s", serverSettings)
	err = http.ListenAndServe(serverSettings, router.Router)
	if err != nil {
		log.Fatalf("Server stating error: %s", err)
	}
}

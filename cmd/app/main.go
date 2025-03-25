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

	serverSettings := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	callbackURL := fmt.Sprintf("http://%s", serverSettings)
	db, err := database.NewDatabase(cfg.DB)
	if err != nil {
		log.Fatalf("DB connection error: %s", err)
	}
	defer db.CloseConnection()

	storage, err := storage.NewStorage(cfg.S3)
	if err != nil {
		log.Fatalf("Storage connection error: %s", err)
	}

	client, err := clients.NewAPIClient(context.Background(), cfg.API, callbackURL)
	if err != nil {
		log.Fatalf("Client connection error: %s", err)
	}

	router := handlers.NewRouter(db, storage, client)
	router.CreateHandlers()

	router.SetReady()
	log.Printf("Server is ready: %s", serverSettings)
	err = http.ListenAndServe(serverSettings, router.Router)
	if err != nil {
		log.Fatalf("Server stating error: %s", err)
	}
}

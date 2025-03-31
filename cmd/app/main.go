package main

import (
	"context"
	"fmt"
	"log"
	"main/internal/clients"
	"main/internal/config"
	"main/internal/database"
	"main/internal/handlers"
	"main/internal/services"
	"main/internal/storage"
	"net/http"

	CORS "github.com/gorilla/handlers"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Config file error: %s", err)
	}

	serverSettings := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	callbackURL := fmt.Sprintf("http://%s", serverSettings)
	db, err := database.NewDatabase(context.Background(), cfg.DB)
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
	service := services.NewService(db, storage)

	router := handlers.NewRouter(service, client)
	router.CreateHandlers()

	cors := CORS.CORS(
		CORS.AllowedOrigins([]string{"127.0.0.1"}),
		CORS.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		CORS.AllowedHeaders([]string{"Content-Type"}),
	)

	router.SetReady()
	log.Printf("Server is ready: %s", serverSettings)
	err = http.ListenAndServe(serverSettings, cors(router.Router))
	if err != nil {
		log.Fatalf("Server stating error: %s", err)
	}
}

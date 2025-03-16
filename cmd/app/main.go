package main

import (
	"fmt"
	"log"
	"main/internal/config"
	"main/internal/database"
	"main/internal/handlers"
	"main/internal/storage"
	"net/http"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Ошибка файла конфигурации: %s", err)
	}

	db, err := database.New(cfg.DB)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %s", err)
	}
	defer db.Close()

	storage, err := storage.NewStorage(cfg.S3)
	if err != nil {
		log.Fatalf("Ошибка подключения к хранилищу: %s", err)
	}

	router := handlers.CreateHandlers(db, storage)

	serverSettings := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	handlers.SetReady()
	log.Printf("Сервер запущен по адресу: %s", serverSettings)
	err = http.ListenAndServe(serverSettings, router)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %s", err)
	}
}

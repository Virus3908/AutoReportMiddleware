package main

import (
	"fmt"
	"log"
	"main/internal/config"
	"main/internal/database"
	"main/internal/handler"
	"net/http"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Ошибка файла конфигурации: %s", err)
	}

	db, _ := database.New(cfg.DBConfig)
	defer db.Close()

	serverSettings := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Сервер запущен по адресу: %s", serverSettings)
	err = http.ListenAndServe(serverSettings, nil)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %s", err)
	}
}

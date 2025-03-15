package main

import (
	"fmt"
	"log"
	"main/internal/config"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Ошибка файла конфигурации: %s", err)
	}
	fmt.Print(cfg.DBConfig.Host)
}
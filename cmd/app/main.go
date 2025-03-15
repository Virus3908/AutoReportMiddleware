package main

import (
	"fmt"
	"log"
	"main/internal/config"

	"database/sql"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Ошибка файла конфигурации: %s", err)
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.User,
		cfg.DBConfig.Password,
		cfg.DBConfig.Database)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")
}

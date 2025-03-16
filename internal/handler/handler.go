package handler

import (
	_ "log"
	"context"
	"main/internal/database"
	"net/http"
	"strings"
	"encoding/json"
	"strconv"
	"github.com/jackc/pgx/v5"
)

func DatabaseHandler(w http.ResponseWriter, r *http.Request, db *database.DB) {
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/database")
	switch r.Method {
	case http.MethodGet:
		getDatabaseHandler(w, r, db)
	case http.MethodPost:
		postDatabaseHandler(w, r, db)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func getDatabaseHandler(w http.ResponseWriter, r *http.Request, db *database.DB) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	var name string
	err = db.Database.QueryRow(context.Background(),
		"SELECT name FROM users WHERE id=$1", userID).Scan(&name)

	switch err {
	case nil:
		response := map[string]string{"id": idStr, "name": name}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	case pgx.ErrNoRows:
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
	default:
		http.Error(w, "Ошибка запроса к базе данных", http.StatusInternalServerError)
	}
}

func postDatabaseHandler(w http.ResponseWriter, r *http.Request, db *database.DB) {
	_, err := db.Database.Exec(context.Background(), "INSERT INTO users (name) VALUES ($1)", r.URL.Path)
	if err != nil {
		http.Error(w, "Ошибка запроса к базе данных", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

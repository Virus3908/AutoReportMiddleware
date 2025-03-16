package handlers

import (
	"main/internal/database"
	"net/http"
	"github.com/gorilla/mux"
	"main/internal/logging"
)

func CreateHandlers(db *database.DataBase) *mux.Router {
	router := mux.NewRouter()
	router.Use(logging.LoggingMidleware)

	router.HandleFunc("/liveness", LivenessHandler)
	router.HandleFunc("/readiness", ReadinessHandler)
	router.HandleFunc("/info", InfoHandler)
	router.HandleFunc("/conversations", func(w http.ResponseWriter, r *http.Request) {
		conversationHandlers(w, r, db)
	})
	router.HandleFunc("/conversations/{id}", func(w http.ResponseWriter, r *http.Request) {
		conversationHandlersWithID(w, r, db)
	})
	return router
}
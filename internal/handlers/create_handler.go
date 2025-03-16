package handlers

import (
	"main/internal/database"
	"main/internal/logging"
	"net/http"
	"github.com/gorilla/mux"
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
	router.HandleFunc("/participant", func(w http.ResponseWriter, r *http.Request) {
		participantHandlers(w, r, db)
	})
	router.HandleFunc("/participant/{id}", func(w http.ResponseWriter, r *http.Request) {
		participantHandlersWithID(w, r, db)
	})
	return router
}
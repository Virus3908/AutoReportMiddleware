package handlers

import (
	"main/internal/database"
	"main/internal/logging"
	"main/internal/storage"
	"main/internal/services"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	Router *mux.Router
	Service *services.Service
}

func NewRouter() *mux.Router {
	return mux.NewRouter()
}

func (r *Router) CreateHandlers(db *database.DataBase, storage *storage.S3Client) {
	r.createLogAndInfoHandlers()
	r.createConversationHandler(db, storage)
	r.createParticipantsHandler(db)

}

func (r *Router) createLogAndInfoHandlers() {
	r.Router.Use(logging.LoggingMidleware)

	r.Router.HandleFunc("/liveness", LivenessHandler).Methods(http.MethodGet)
	r.Router.HandleFunc("/readiness", ReadinessHandler).Methods(http.MethodGet)
	r.Router.HandleFunc("/info", InfoHandler).Methods(http.MethodGet)
}

func (r *Router) createConversationHandler(db *database.DataBase, storage *storage.S3Client) {
	r.Router.HandleFunc("/conversations", func(w http.ResponseWriter, r *http.Request) {
		getConversationsHandler(w, r, db)
	}).Methods(http.MethodGet)
	r.Router.HandleFunc("/conversations", func(w http.ResponseWriter, r *http.Request) {
		createConversationHandler(w, r, db, storage)
	}).Methods(http.MethodPost)
}

func (r *Router) createParticipantsHandler(db *database.DataBase) {
	r.Router.HandleFunc("/conversations/{id}", func(w http.ResponseWriter, r *http.Request) {
		conversationHandlersWithID(w, r, db)
	})
	r.Router.HandleFunc("/participant", func(w http.ResponseWriter, r *http.Request) {
		participantHandlers(w, r, db)
	})
	r.Router.HandleFunc("/participant/{id}", func(w http.ResponseWriter, r *http.Request) {
		participantHandlersWithID(w, r, db)
	})
}
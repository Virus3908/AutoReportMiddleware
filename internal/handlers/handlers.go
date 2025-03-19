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
	DB database.Database
	Storage storage.Storage
	Service services.Service
}

func NewRouter(db database.Database, storage storage.Storage) *Router {
	return &Router{
		Router: mux.NewRouter(),
		Service: nil,
		Storage: storage,
		DB: db,
	}
}

func (r *Router) CreateHandlers() {
	r.createLogAndInfoHandlers()
	r.conversationHandler()
	r.participantsHandler()

}

func (r *Router) createLogAndInfoHandlers() {
	r.Router.Use(logging.LoggingMidleware)

	r.Router.HandleFunc("/liveness", LivenessHandler).Methods(http.MethodGet)
	r.Router.HandleFunc("/readiness", ReadinessHandler).Methods(http.MethodGet)
	r.Router.HandleFunc("/info", InfoHandler).Methods(http.MethodGet)
}

func (router *Router) conversationHandler() {
	router.Router.HandleFunc("/conversations", func(w http.ResponseWriter, r *http.Request) {
		router.getConversationsHandler(w, r)
	}).Methods(http.MethodGet)
	router.Router.HandleFunc("/conversations", func(w http.ResponseWriter, r *http.Request) {
		router.createConversationHandler(w, r)
	}).Methods(http.MethodPost)
}

func (router *Router) participantsHandler() {
	router.Router.HandleFunc("/conversations/{id}", func(w http.ResponseWriter, r *http.Request) {
		router.conversationHandlersWithID(w, r)
	})
	router.Router.HandleFunc("/participant", func(w http.ResponseWriter, r *http.Request) {
		router.participantHandlers(w, r)
	})
	router.Router.HandleFunc("/participant/{id}", func(w http.ResponseWriter, r *http.Request) {
		router.participantHandlersWithID(w, r)
	})
}
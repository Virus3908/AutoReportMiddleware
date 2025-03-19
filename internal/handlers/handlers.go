package handlers

import (
	"main/internal/database"
	"main/internal/logging"
	"main/internal/services"
	"main/internal/storage"

	"net/http"
	"sync/atomic"

	"github.com/gorilla/mux"
)

var ready int32

type Router interface {
	CreateHandlers()
	SetReady()
	// createLogAndInfoHandlers()
	// conversationHandler()
	// participantsHandler()
}

type RouterStruct struct {
	Router  *mux.Router
	DB      database.Database
	Storage storage.Storage
	Service services.Service
}

func NewRouter(db database.Database, storage storage.Storage) *RouterStruct {
	return &RouterStruct{
		Router:  mux.NewRouter(),
		Service: nil,
		Storage: storage,
		DB:      db,
	}
}

func (_ *RouterStruct) SetReady() {
	atomic.StoreInt32(&ready, 1)
}

func (r *RouterStruct) CreateHandlers() {
	r.logAndInfoHandlers()
	r.conversationHandlers()
	r.participantsHandlers()
	r.promtHandlers()
}

func (r *RouterStruct) logAndInfoHandlers() {
	r.Router.Use(logging.LoggingMidleware)

	r.Router.HandleFunc("/liveness", r.livenessHandler).Methods(http.MethodGet)
	r.Router.HandleFunc("/readiness", r.readinessHandler).Methods(http.MethodGet)
	r.Router.HandleFunc("/info", r.infoHandler).Methods(http.MethodGet)
}

func (r *RouterStruct) conversationHandlers() {
	r.Router.HandleFunc("/api/conversations", r.getConversationsHandler).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/conversations", r.createConversationHandler).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/conversations/{id}", r.getConversationByIDHandler).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/conversations/{id}", r.updateConversationNameByIDHandler).Methods(http.MethodPut)
	r.Router.HandleFunc("/api/conversations/{id}", r.deleteConversationByIDHandler).Methods(http.MethodDelete)
}

func (r *RouterStruct) participantsHandlers() {
	r.Router.HandleFunc("/api/participant", r.getParticipantsHandler).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/participant/{id}", r.getParticipantByIDHandler).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/participant", r.createParticipantHandler).Methods(http.MethodPost)	
	r.Router.HandleFunc("/api/participant/{id}", r.updateParticipantByIDHandler).Methods(http.MethodPut)
	r.Router.HandleFunc("/api/participant", r.deleteParticipantByIDHandler).Methods(http.MethodDelete)
}

func (r *RouterStruct) promtHandlers() {
	r.Router.HandleFunc("/api/promt", r.getPromtsHandler).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/promt/{id}", r.getPromtByIDHandler).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/promt", r.createPromtHandler).Methods(http.MethodPost)	
	r.Router.HandleFunc("/api/promt/{id}", r.updatePromtByIDHandler).Methods(http.MethodPut)
	r.Router.HandleFunc("/api/promt", r.deletePromtByIDHandler).Methods(http.MethodDelete)
}
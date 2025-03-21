package handlers

import (
	"main/internal/database"
	"main/internal/logging"
	"main/internal/storage"
	"main/internal/clients"

	"net/http"
	"sync/atomic"

	"github.com/gorilla/mux"
)

type Router interface {
	CreateHandlers()
	SetReady()
	// createLogAndInfoHandlers()
	// conversationHandler()
	// participantsHandler()
}

type RouterStruct struct {
	Client 	clients.Client
	Router  *mux.Router
	DB      database.Database
	Storage storage.Storage
	CallbackURL string
	ready int32
}

func NewRouter(db database.Database, storage storage.Storage, client clients.Client) *RouterStruct {
	return &RouterStruct{
		Client: client,
		Router:  mux.NewRouter(),
		Storage: storage,
		DB:      db,
		ready: 0,
	}
}

func (r *RouterStruct) SetReady() {
	atomic.StoreInt32(&r.ready, 1)
}

func (r *RouterStruct) CreateHandlers() {
	r.logAndInfoHandlers()
	r.conversationHandlers()
	r.participantsHandlers()
	r.promtHandlers()
	r.updateTaskHandlers()
	r.createTaskHandlers()
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

func (r *RouterStruct) updateTaskHandlers() {
	r.Router.HandleFunc("/api/update/convert/{id}", r.acceptConvertFileHandler).Methods(http.MethodPut)
}

func (r *RouterStruct) createTaskHandlers() {
	r.Router.HandleFunc("/api/create/convert/{id}", r.createConvertFileTaskHandlerd).Methods(http.MethodPost)
}
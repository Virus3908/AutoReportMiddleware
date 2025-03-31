package handlers

import (
	"main/internal/clients"
	"main/internal/logging"
	"main/internal/services"

	"net/http"
	"sync/atomic"

	"github.com/gorilla/mux"
)

type Router interface {
	CreateHandlers()
	SetReady()
}

type RouterStruct struct {
	Client  clients.Client
	Router  *mux.Router
	Service *services.ServicesStruct
	ready   int32
}

func NewRouter(service *services.ServicesStruct, client clients.Client) *RouterStruct {
	return &RouterStruct{
		Client:  client,
		Router:  mux.NewRouter(),
		Service: service,
		ready:   0,
	}
}

func (r *RouterStruct) SetReady() {
	atomic.StoreInt32(&r.ready, 1)
}

func (r *RouterStruct) CreateHandlers() {
	r.logAndInfoHandlers()
	r.participantsHandlers()
	r.promptsHandlers()
	r.conversationsHandlers()
}

func (r *RouterStruct) logAndInfoHandlers() {
	r.Router.Use(logging.LoggingMidleware)

	r.Router.HandleFunc("/liveness", r.livenessHandler).Methods(http.MethodGet)
	r.Router.HandleFunc("/readiness", r.readinessHandler).Methods(http.MethodGet)
	r.Router.HandleFunc("/info", r.infoHandler).Methods(http.MethodGet)
}

func (r *RouterStruct) participantsHandlers() {
	r.Router.HandleFunc("/api/participants",
		wrapperGetHandler(r.Service.Participant.GetAll),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/participants/{id}",
		wrapperGetByIDHandler(r.Service.Participant.GetByID),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/participants",
		wrapperCreateHandler(r.Service.Participant.Create),
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/participants/{id}",
		wrapperUpdateHandler(r.Service.Participant.Update),
	).Methods(http.MethodPut)
	r.Router.HandleFunc("/api/participants/{id}",
		wrapperDeleteHandler(r.Service.Participant.Delete),
	).Methods(http.MethodDelete)
}

func (r *RouterStruct) promptsHandlers() {
	r.Router.HandleFunc("/api/prompts",
		wrapperGetHandler(r.Service.Prompt.GetAll),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/prompts/{id}",
		wrapperGetByIDHandler(r.Service.Prompt.GetByID),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/prompts",
		wrapperCreateHandler(r.Service.Prompt.Create),
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/prompts/{id}",
		wrapperUpdateHandler(r.Service.Prompt.Update),
	).Methods(http.MethodPut)
	r.Router.HandleFunc("/api/prompts/{id}",
		wrapperDeleteHandler(r.Service.Prompt.Delete),
	).Methods(http.MethodDelete)
}

func (r *RouterStruct) conversationsHandlers() {
	r.Router.HandleFunc("/api/conversations",
		wrapperGetHandler(r.Service.Conversation.GetAll),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/conversations/{id}",
		wrapperGetByIDHandler(r.Service.Conversation.GetByID),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/conversations",
		r.createConversationHandler,
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/conversations/{id}",
		wrapperUpdateHandler(r.Service.Conversation.Update),
	).Methods(http.MethodPut)
	r.Router.HandleFunc("/api/conversations/{id}",
		wrapperDeleteHandler(r.Service.Conversation.Delete),
	).Methods(http.MethodDelete)
}
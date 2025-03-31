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
	r.promtsHandlers()
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
		wrapperGetHandler(r.Service.Participant.GetParticipants),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/participants/{id}",
		wrapperGetByIDHandler(r.Service.Participant.GetParticipantByID),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/participants",
		wrapperCreateHandler(r.Service.Participant.CreateParticipant),
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/participants/{id}",
		wrapperUpdateHandler(r.Service.Participant.UpdateParticipant),
	).Methods(http.MethodPut)
	r.Router.HandleFunc("/api/participants/{id}",
		wrapperDeleteHandler(r.Service.Participant.DeleteParticipant),
	).Methods(http.MethodDelete)
}

func (r *RouterStruct) promtsHandlers() {
	r.Router.HandleFunc("/api/promts",
		wrapperGetHandler(r.Service.Promt.GetPromts),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/promts/{id}",
		wrapperGetByIDHandler(r.Service.Promt.GetPromtByID),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/promts",
		wrapperCreateHandler(r.Service.Promt.CreatePromt),
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/promts/{id}",
		wrapperUpdateHandler(r.Service.Promt.UpdatePromt),
	).Methods(http.MethodPut)
	r.Router.HandleFunc("/api/promts/{id}",
		wrapperDeleteHandler(r.Service.Promt.DeletePromt),
	).Methods(http.MethodDelete)
}

func (r *RouterStruct) conversationsHandlers() {
	r.Router.HandleFunc("/api/conversations",
		wrapperGetHandler(r.Service.Conversation.GetConversations),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/conversations/{id}",
		wrapperGetByIDHandler(r.Service.Conversation.GetConversationByID),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/conversations",
		r.createConversationHandler,
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/conversations/{id}",
		wrapperUpdateHandler(r.Service.Conversation.UpdateConversation),
	).Methods(http.MethodPut)
	r.Router.HandleFunc("/api/conversations/{id}",
		wrapperDeleteHandler(r.Service.Conversation.DeleteConversation),
	).Methods(http.MethodDelete)
}

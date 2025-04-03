package handlers

import (
	"main/internal/logging"
	"main/internal/services"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/mux"
)

type Router interface {
	CreateHandlers()
	SetReady()
}

type RouterStruct struct {
	Router  *mux.Router
	Service *services.ServicesStruct
	ready   int32
}

func NewRouter(service *services.ServicesStruct) *RouterStruct {
	return &RouterStruct{
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
		wrapperGetHandler(r.Service.CrudService.Participant.GetAll),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/participants/{id}",
		wrapperGetByIDHandler(r.Service.CrudService.Participant.GetByID),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/participants",
		wrapperCreateHandler(r.Service.CrudService.Participant.Create),
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/participants/{id}",
		wrapperUpdateHandler(r.Service.CrudService.Participant.Update),
	).Methods(http.MethodPut)
	r.Router.HandleFunc("/api/participants/{id}",
		wrapperDeleteHandler(r.Service.CrudService.Participant.Delete),
	).Methods(http.MethodDelete)
}

func (r *RouterStruct) promptsHandlers() {
	r.Router.HandleFunc("/api/prompts",
		wrapperGetHandler(r.Service.CrudService.Prompt.GetAll),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/prompts/{id}",
		wrapperGetByIDHandler(r.Service.CrudService.Prompt.GetByID),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/prompts",
		wrapperCreateHandler(r.Service.CrudService.Prompt.Create),
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/prompts/{id}",
		wrapperUpdateHandler(r.Service.CrudService.Prompt.Update),
	).Methods(http.MethodPut)
	r.Router.HandleFunc("/api/prompts/{id}",
		wrapperDeleteHandler(r.Service.CrudService.Prompt.Delete),
	).Methods(http.MethodDelete)
}

func (r *RouterStruct) conversationsHandlers() {
	r.Router.HandleFunc("/api/conversations",
		wrapperGetHandler(r.Service.CrudService.Conversation.GetAll),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/conversations/{id}",
		wrapperGetByIDHandler(r.Service.CrudService.Conversation.GetByID),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/conversations",
		r.createConversationHandler,
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/conversations/{id}",
		wrapperUpdateHandler(r.Service.CrudService.Conversation.Update),
	).Methods(http.MethodPut)
	r.Router.HandleFunc("/api/conversations/{id}",
		wrapperDeleteHandler(r.Service.ConversationsService.DeleteConversations),
	).Methods(http.MethodDelete)
}

func respondWithError(w http.ResponseWriter, msg string, err error, status int) {
    log.Printf("[ERROR] %s: %v", msg, err)
    http.Error(w, msg, status)
}
package handlers

import (
	"log"
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
	r.taskHandlers()
}

func (r *RouterStruct) logAndInfoHandlers() {
	r.Router.Use(logging.LoggingMidleware)

	r.Router.HandleFunc("/liveness", r.livenessHandler).Methods(http.MethodGet)
	r.Router.HandleFunc("/readiness", r.readinessHandler).Methods(http.MethodGet)
	r.Router.HandleFunc("/info", r.infoHandler).Methods(http.MethodGet)
}

func (r *RouterStruct) participantsHandlers() {
	r.Router.HandleFunc("/api/participants",
		wrapperReturningData(r.Service.CrudService.Participant.GetAll),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/participants/{id}",
		wrapperWithIDReturningData(r.Service.CrudService.Participant.GetByID),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/participants",
		wrapperWithPayload(r.Service.CrudService.Participant.Create),
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/participants/{id}",
		wrapperWithIDAndPayload(r.Service.CrudService.Participant.Update),
	).Methods(http.MethodPut)
	r.Router.HandleFunc("/api/participants/{id}",
		wrapperWithID(r.Service.CrudService.Participant.Delete),
	).Methods(http.MethodDelete)
}

func (r *RouterStruct) promptsHandlers() {
	r.Router.HandleFunc("/api/prompts",
		wrapperReturningData(r.Service.CrudService.Prompt.GetAll),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/prompts/{id}",
		wrapperWithIDReturningData(r.Service.CrudService.Prompt.GetByID),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/prompts",
		wrapperWithPayload(r.Service.CrudService.Prompt.Create),
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/prompts/{id}",
		wrapperWithIDAndPayload(r.Service.CrudService.Prompt.Update),
	).Methods(http.MethodPut)
	r.Router.HandleFunc("/api/prompts/{id}",
		wrapperWithID(r.Service.CrudService.Prompt.Delete),
	).Methods(http.MethodDelete)
}

func (r *RouterStruct) conversationsHandlers() {
	r.Router.HandleFunc("/api/conversations",
		wrapperReturningData(r.Service.CrudService.Conversation.GetAll),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/conversations/{id}",
		wrapperWithIDReturningData(r.Service.CrudService.Conversation.GetByID),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/conversations",
		r.createConversationHandler,
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/conversations/{id}",
		wrapperWithIDAndPayload(r.Service.CrudService.Conversation.Update),
	).Methods(http.MethodPut)
	r.Router.HandleFunc("/api/conversations/{id}",
		wrapperWithID(r.Service.ConversationsService.DeleteConversations),
	).Methods(http.MethodDelete)
}

func (r *RouterStruct) taskHandlers() {
	r.Router.HandleFunc("/api/task/create/convert/{id}",
		wrapperWithID(r.Service.TaskService.CreateConvertTask),
	).Methods(http.MethodPost)
}

func respondWithError(w http.ResponseWriter, msg string, err error, status int) {
	log.Printf("[ERROR] %s: %v", msg, err)
	http.Error(w, msg, status)
}

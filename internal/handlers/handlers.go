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
	r.participantHandlers()
}

func (r *RouterStruct) logAndInfoHandlers() {
	r.Router.Use(logging.LoggingMidleware)

	r.Router.HandleFunc("/liveness", r.livenessHandler).Methods(http.MethodGet)
	r.Router.HandleFunc("/readiness", r.readinessHandler).Methods(http.MethodGet)
	r.Router.HandleFunc("/info", r.infoHandler).Methods(http.MethodGet)
}

func (r *RouterStruct) participantHandlers() {
	r.Router.HandleFunc("/api/participants",
		simpleGetHandler(r.Service.Participant.Get),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/participants/{id}",
		simpleGetByIDHandler(r.Service.Participant.GetByID),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/participants",
		simpleCreateHandler(r.Service.Participant.Create),
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/participants/{id}",
		simpleUpdateHandler(r.Service.Participant.Update),
	).Methods(http.MethodPut)
	r.Router.HandleFunc("/api/participants/{id}",
		simpleDeleteHandler(r.Service.Participant.Delete),
	).Methods(http.MethodDelete)
}

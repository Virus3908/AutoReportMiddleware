package handlers

import (
	"log"
	"main/internal/services"
	"net/http"

	"github.com/gorilla/mux"
)

type RouterStruct struct {
	Router  *mux.Router
	Service *services.ServiceStruct
}

func New(service *services.ServiceStruct, mwf []mux.MiddlewareFunc) *RouterStruct {
	router := RouterStruct{
		Router: mux.NewRouter(),
		Service: service,
	}
	router.Router.Use(mwf...)
	router.createHandlers()
	return &router
}

func (router *RouterStruct) GetRouter() *mux.Router {
	return router.Router
}

func (r *RouterStruct) createHandlers() {
	r.participantsHandlers()
	r.promptsHandlers()
	r.conversationsHandlers()
	r.taskHandlers()
}

func (r *RouterStruct) participantsHandlers() {

}

func (r *RouterStruct) promptsHandlers() {

}

func (r *RouterStruct) conversationsHandlers() {
	r.Router.HandleFunc("/api/conversations",
		wrapperReturningData(r.Service.Repo.GetConversations),
	).Methods(http.MethodGet)
// 	r.Router.HandleFunc("/api/conversations/{id}",
// 		wrapperWithIDReturningData(r.Service.CrudService.Conversation.GetByID),
// 	).Methods(http.MethodGet)
// 	r.Router.HandleFunc("/api/conversations",
// 		r.createConversationHandler,
// 	).Methods(http.MethodPost)
// 	r.Router.HandleFunc("/api/conversations/{id}",
// 		wrapperWithIDAndPayload(r.Service.CrudService.Conversation.Update),
// 	).Methods(http.MethodPut)
// 	r.Router.HandleFunc("/api/conversations/{id}",
// 		wrapperWithID(r.Service.ConversationsService.DeleteConversations),
// 	).Methods(http.MethodDelete)
}

func (r *RouterStruct) taskHandlers() {
	// r.Router.HandleFunc("/api/task/create/convert/{id}",
	// 	wrapperWithID(r.Service.TaskService.CreateConvertTask),
	// ).Methods(http.MethodPost)
}

func respondWithError(w http.ResponseWriter, msg string, err error, status int) {
	log.Printf("[ERROR] %s: %v", msg, err)
	http.Error(w, msg, status)
}

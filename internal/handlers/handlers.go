package handlers

import (
	"main/internal/common/interfaces"
	"main/internal/logger"
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
		Router:  mux.NewRouter(),
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
	r.segmentHandlers()
	r.transcriptionHandlers()
}

func (r *RouterStruct) participantsHandlers() {
	r.Router.HandleFunc("/api/participants",
		wrapperWithPayload(r.Service.Participants.CreateParticipant),
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/participants",
		wrapperReturningData(r.Service.Participants.GetParticipants),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/participants/{id}",
		wrapperWithID(r.Service.Participants.DeleteParticipantByID),
	).Methods(http.MethodDelete)
	r.Router.HandleFunc("/api/participants/{id}",
		wrapperWithIDAndPayload(r.Service.Participants.UpdateParticipantByID),
	).Methods(http.MethodPatch)
	r.Router.HandleFunc("/api/participants/{id}",
		wrapperWithIDReturningData(r.Service.Participants.GetParticipantByID),
	).Methods(http.MethodGet)
}

func (r *RouterStruct) segmentHandlers() {
	r.Router.HandleFunc("/api/segments/{id}",
		wrapperWithIDAndPayload(r.Service.Conversations.AssignParticipantToSegment),
	).Methods(http.MethodPatch)
}

func (r *RouterStruct) promptsHandlers() {
	r.Router.HandleFunc("/api/prompts",
		wrapperReturningData(r.Service.Prompts.GetPrompts),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/prompts",
		wrapperWithPayload(r.Service.Prompts.CreatePrompt),
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/prompts/{id}",
		wrapperWithID(r.Service.Prompts.DeletePromptByID),
	).Methods(http.MethodDelete)
	r.Router.HandleFunc("/api/prompts/{id}",
		wrapperWithIDAndPayload(r.Service.Prompts.UpdatePromptByID),
	).Methods(http.MethodPatch)
}

func (r *RouterStruct) conversationsHandlers() {
	r.Router.HandleFunc("/api/conversations",
		wrapperReturningData(r.Service.Conversations.GetConversations),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/conversations",
		r.CreateConversation,
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/conversations/{id}",
		wrapperWithIDReturningData(r.Service.Conversations.GetConversationDetails),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/conversations/{id}",
		wrapperWithID(r.Service.Conversations.DeleteConversationByID),
	).Methods(http.MethodDelete)
	r.Router.HandleFunc("/api/conversations/{id}",
		wrapperWithIDAndPayload(r.Service.Conversations.UpdateConversationNameByID),
	).Methods(http.MethodPatch)
}

func (r *RouterStruct) taskHandlers() {
	r.Router.HandleFunc("/api/task/create/convert/{id}",
		wrapperWithID(r.Service.Tasks.CreateConvertTask),
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/task/create/diarize/{id}",
		wrapperWithID(r.Service.Tasks.CreateDiarizeTask),
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/task/create/transcription/{id}",
		wrapperWithID(r.Service.Tasks.CreateTranscribeTask),
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/task/create/semireport/{id}",
		wrapperWithIDAndPayload(r.Service.Tasks.CreateSemiReportTask),
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/task/create/report/{id}",
		wrapperWithIDAndPayload(r.Service.Tasks.CreateReportTask),
	).Methods(http.MethodPost)
}

func (r *RouterStruct) transcriptionHandlers() {
	r.Router.HandleFunc("/api/transcription/update/{id}",
		wrapperWithIDAndPayload(r.Service.Conversations.UpdateTranscriptionTextByID),
	).Methods(http.MethodPatch)
}

func respondWithError(w http.ResponseWriter, r *http.Request, err error, status int) {
	logger := logger.GetLoggerFromContext(r.Context())
	logger.Error("error during request", interfaces.LogField{Key: "error", Value: err.Error()})
	http.Error(w, err.Error(), status)
}

func (h *RouterStruct) CreateConversation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(200 << 20) // swap to var
	if err != nil {
		respondWithError(w, r, err, http.StatusBadRequest)
		return
	}

	conversationsName := r.FormValue("conversation_name")

	file, header, err := r.FormFile("file")
	if err != nil {
		respondWithError(w, r, err, http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = h.Service.Conversations.CreateConversation(r.Context(), conversationsName, header.Filename, file)
	if err != nil {
		respondWithError(w, r, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

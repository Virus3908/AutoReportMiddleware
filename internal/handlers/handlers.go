package handlers

import (
	"log"
	"main/internal/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
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
	r.callbackHandlers()
}

func (r *RouterStruct) participantsHandlers() {

}

func (r *RouterStruct) promptsHandlers() {

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
		wrapperWithID(r.Service.Conversations.DeleteConversation),
	).Methods(http.MethodDelete)
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
}

func (r *RouterStruct) callbackHandlers() {
	r.Router.HandleFunc("/api/task/update/convert/{id}",
		r.UpdateConvert,
	).Methods(http.MethodPatch)
	r.Router.HandleFunc("/api/task/update/diarize/{id}",
		wrapperWithIDAndPayload(r.Service.TaskCallbackReceiver.HandleDiarizeCallback),
	).Methods(http.MethodPatch)
	r.Router.HandleFunc("/api/task/update/transcription/{id}",
		wrapperWithIDAndPayload(r.Service.TaskCallbackReceiver.HandleTransctiprionCallback),
	).Methods(http.MethodPatch)
}

func respondWithError(w http.ResponseWriter, msg string, err error, status int) {
	log.Printf("[ERROR] %s: %v", msg, err)
	http.Error(w, msg, status)
}

func (h *RouterStruct) CreateConversation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(200 << 20)
	if err != nil {
		respondWithError(w, "can't parse form", err, http.StatusBadRequest)
		return
	}

	conversationsName := r.FormValue("conversation_name")

	file, header, err := r.FormFile("file")
	if err != nil {
		respondWithError(w, "file not found", err, http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = h.Service.Conversations.CreateConversation(r.Context(), conversationsName, header.Filename, file)
	if err != nil {
		respondWithError(w, "failed to create conversation", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *RouterStruct) UpdateConvert(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(200 << 20)
	if err != nil {
		respondWithError(w, "can't parse form", err, http.StatusBadRequest)
		return
	}
	audioLenStr := r.FormValue("audio_len")
	audioLen, err := strconv.ParseFloat(strings.TrimSpace(audioLenStr), 64)
	if err != nil {
		respondWithError(w, "can't parse audio len", err, http.StatusBadRequest)
		return
	}

	strTaskID := mux.Vars(r)["id"]
	taskID, err := uuid.Parse(strTaskID)
	if err != nil {
		respondWithError(w, "invalid task ID", err, http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		respondWithError(w, "file not found", err, http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = h.Service.TaskCallbackReceiver.HandleConvertCallback(r.Context(), taskID, file, header.Filename, audioLen)
	if err != nil {
		respondWithError(w, "failed to update convert", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

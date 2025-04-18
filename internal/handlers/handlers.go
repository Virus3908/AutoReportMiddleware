package handlers

import (
	"log"
	"main/internal/services"
	messages "main/pkg/messages/proto"
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
	r.segmentHandlers()
	r.callbackHandlers()
}

func (r *RouterStruct) participantsHandlers() {
	r.Router.HandleFunc("/api/participants",
		wrapperWithPayload(r.Service.Conversations.CreateParticipant),
	).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/participants",
		wrapperReturningData(r.Service.Conversations.GetParticipants),
	).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/participants/{id}",
		wrapperWithID(r.Service.Conversations.DeleteParticipantByID),
	).Methods(http.MethodDelete)
}

func (r *RouterStruct) segmentHandlers() {
	r.Router.HandleFunc("/api/segments/{id}",
		wrapperWithIDAndPayload(r.Service.Conversations.AssignParticipantToSegment),
	).Methods(http.MethodPatch)
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
		wrapperWithID(r.Service.Conversations.DeleteConversationByID),
	).Methods(http.MethodDelete)
	r.Router.HandleFunc("/api/transcription/update/{id}",
		wrapperWithIDAndPayload(r.Service.Conversations.UpdateTranscriptionTextByID),	
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
}

func (r *RouterStruct) callbackHandlers() {
	r.Router.HandleFunc("/api/task/update/convert/{id}",
		r.UpdateConvert,
	).Methods(http.MethodPatch)
	r.Router.HandleFunc("/api/task/update/diarize/{id}",
		r.handleDiarizeCallback,
	).Methods(http.MethodPatch)
	r.Router.HandleFunc("/api/task/update/transcription/{id}",
		r.handleTranscriptionCallback,
	).Methods(http.MethodPatch)
	r.Router.HandleFunc("/api/task/update/error/{id}",
		r.handleErrorCallback,
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

	err = h.Service.Tasks.HandleConvertCallback(r.Context(), taskID, file, header.Filename, audioLen)
	if err != nil {
		respondWithError(w, "failed to update convert", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (r *RouterStruct) handleDiarizeCallback(w http.ResponseWriter, req *http.Request) {
	handleProtoRequest(w, req, &messages.SegmentsTaskResponse{}, r.Service.Tasks.HandleDiarizeCallback)
}

func (r *RouterStruct) handleTranscriptionCallback(w http.ResponseWriter, req *http.Request) {
	handleProtoRequest(w, req, &messages.TranscriptionTaskResponse{}, r.Service.Tasks.HandleTransctiprionCallback)
}

func (r *RouterStruct) handleErrorCallback(w http.ResponseWriter, req *http.Request) {
	handleProtoRequest(w, req, &messages.ErrorTaskResponse{}, r.Service.Tasks.HandleErrorCallback)
}
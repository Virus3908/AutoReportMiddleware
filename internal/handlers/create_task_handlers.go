package handlers

import (
	"encoding/json"
	"log"
	"main/internal/clients"
	"main/internal/repositories"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (router *RouterStruct) createConvertFileTaskHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	strID, ok := params["id"]
	if !ok {
		http.Error(w, "missing id in request", http.StatusBadRequest)
		return
	}
	conversationID, err := uuid.Parse(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, rollback, commit, err := router.DB.StartTransaction()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rollback()

	fileURL, err := tx.GetConversationFileURL(r.Context(), conversationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	taskUUID, err := router.Client.CreateTaskConvertFileAndGetTaskID(r.Context(), fileURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if taskUUID == nil {
		http.Error(w, "taskUUID is nil", http.StatusInternalServerError)
		return
	}

	createTask := repositories.CreateConvertTaskParams{
		ConversationsID: conversationID,
		TaskID: *taskUUID,
	}
	err = tx.CreateConvertTask(r.Context(), createTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Convert task created",
		"task_id":  taskUUID.String(),
		"file_url": fileURL,
	})
	commit()
}

func (router *RouterStruct) createDiarizeTaskeHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	strID, ok := params["id"]
	if !ok {
		http.Error(w, "missing id in request", http.StatusBadRequest)
		return
	}
	conversationID, err := uuid.Parse(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, rollback, commit, err := router.DB.StartTransaction()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rollback()

	fileURL, err := tx.GetConvertFileURL(r.Context(), conversationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if fileURL == nil {
		http.Error(w, "file is not converted", http.StatusInternalServerError)
		return
	}

	taskUUID, err := router.Client.CreateTaskDiarizeFileAndGetTaskID(r.Context(), *fileURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if taskUUID == nil {
		http.Error(w, "taskUUID is nil", http.StatusInternalServerError)
		return
	}

	createTask := repositories.CreateDiarizeTaskParams{
		ConversationID: conversationID,
		TaskID: *taskUUID,
	}
	err = tx.CreateDiarizeTask(r.Context(), createTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Diarize task created",
		"task_id":  taskUUID.String(),
	})
	commit()
}

func (router *RouterStruct) createTranscribeTaskHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	strID, ok := params["id"]
	if !ok {
		http.Error(w, "missing id in request", http.StatusBadRequest)
		return
	}
	conversationID, err := uuid.Parse(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, rollback, commit, err := router.DB.StartTransaction()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rollback()

	segments, err := tx.GetConversationsSegments(r.Context(), conversationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fileURL, err := tx.GetConvertFileURL(r.Context(), conversationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Print(segments)
	for _, segment := range segments{
		taskID, err := router.Client.CreateTaskTranscribeSegmentFileAndGetTaskID(r.Context(), *fileURL, clients.Segment{
			StartTime: segment.StartTime,
			EndTime: segment.EndTime,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		if taskID == nil {
			http.Error(w, "taskUUID is nil", http.StatusInternalServerError)
			return
		}
		createTask := repositories.CreateTranscribeTaskeParams{
			ConversationID: conversationID,
			SegmentID: segment.ID,
			TaskID: *taskID,
		}
		err = tx.CreateTranscribeTaske(r.Context(), createTask)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Transcribe task created",
	})
	commit()
}
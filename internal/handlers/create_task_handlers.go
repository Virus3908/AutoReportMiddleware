package handlers

import (
	"main/internal/repositories"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"encoding/json"
)

func (router *RouterStruct) createConvertFileTaskHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	strID := params["id"]

	UUID, err := uuid.Parse(strID)
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

	fileURL, err := tx.GetConversationFileURL(r.Context(), UUID)
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
		ConversationsID: UUID,
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
	strID := params["id"]

	UUID, err := uuid.Parse(strID)
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

	fileURL, err := tx.GetConvertFileURL(r.Context(), UUID)
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
		ConversationID: UUID,
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
		"file_url": fileURL,
	})
	commit()
}

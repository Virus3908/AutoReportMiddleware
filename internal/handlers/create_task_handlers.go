package handlers

import (
	"context"
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

	fileURL, err := tx.GetConversationFileURL(context.Background(), UUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	taskUUID, err := router.Client.CreateTaskConvertFileAndGetTaskID(context.Background(), fileURL)
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
	err = tx.CreateConvertTask(context.Background(), createTask)
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

	fileURL, err := tx.GetConvertFileURL(context.Background(), UUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if fileURL == nil {
		http.Error(w, "file is not converted", http.StatusInternalServerError)
		return
	}

	taskUUID, err := router.Client.CreateTaskDiarizeFileAndGetTaskID(context.Background(), *fileURL)
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
	err = tx.CreateDiarizeTask(context.Background(), createTask)
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

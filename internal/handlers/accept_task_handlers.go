package handlers

import (
	"io"
	"main/internal/repositories"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
)

func (router *RouterStruct) acceptConvertFileHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	strID, ok := params["id"]
	if !ok {
		http.Error(w, "missing id in request", http.StatusBadRequest)
		return
	}
	taskID, err := uuid.Parse(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1048576)) // 1MB max
	if err != nil {
		http.Error(w, "Error request reading: " + err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	
	fileURL, audioLen, err := router.Client.GetConvertedFileURLAudioLen(body)
	if err != nil {
		http.Error(w, "Error reading body request: " + err.Error(), http.StatusBadRequest)
		return
	}

	tx, rollback, commit, err := router.DB.StartTransaction()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rollback()

	updatedData := repositories.UpdateConvertTaskParams{
		FileUrl: fileURL,
		TaskID: taskID,
		AudioLen: audioLen,
	}

	err = tx.UpdateConvertTask(r.Context(), updatedData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
}

func (router *RouterStruct) acceptDiarizeSegmentsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	strID, ok := params["id"]
	if !ok {
		http.Error(w, "missing id in request", http.StatusBadRequest)
		return
	}
	taskID, err := uuid.Parse(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1048576)) // 1MB max
	if err != nil {
		http.Error(w, "Error request reading: " + err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	segments, err := router.Client.GetDiarizationSegments(body)
	if err != nil {
		http.Error(w, "Error reading body request: " + err.Error(), http.StatusBadRequest)
		return
	}

	tx, rollback, commit, err := router.DB.StartTransaction()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rollback()
	
	conversationID, err := tx.GetConversationIDByDiarizeTaskID(r.Context(), taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, segment := range segments {
		err = tx.CreateSegments(r.Context(), repositories.CreateSegmentsParams{
			ConversationID: conversationID,
			StartTime:      segment.StartTime,
			EndTime:        segment.EndTime,
			Speaker:      segment.Speaker,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	commit()
}

func (router *RouterStruct) acceptTranscibeHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	strID, ok := params["id"]
	if !ok {
		http.Error(w, "missing id in request", http.StatusBadRequest)
		return
	}
	taskID, err := uuid.Parse(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1048576)) // 1MB max
	if err != nil {
		http.Error(w, "Error request reading: " + err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	
	transcription, err := router.Client.GetMessage(body)
	if err != nil {
		http.Error(w, "Error reading body request: " + err.Error(), http.StatusBadRequest)
		return
	}

	tx, rollback, commit, err := router.DB.StartTransaction()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rollback()

	updatedData := repositories.UpdateTranscribeByTaskIDParams{
		Transcription: transcription,
		TaskID: taskID,
	}

	err = tx.UpdateTranscribeByTaskID(r.Context(), updatedData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
}
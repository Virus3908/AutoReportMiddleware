package handlers

import (
	"io"
	"context"
	"main/internal/repositories"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
)

func (router *RouterStruct) acceptConvertFileHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	strID := params["id"]

	UUID, err := uuid.Parse(strID)
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
		TaskID: UUID,
		AudioLen: audioLen,
	}

	err = tx.UpdateConvertTask(context.Background(), updatedData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
}
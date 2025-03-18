package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"main/internal/common"
	"main/internal/database"
	"main/internal/repositories"
	"main/internal/storage"
	"net/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func conversationHandlers(w http.ResponseWriter, r *http.Request, db *database.DataBase, storage *storage.S3Client) {
	switch r.Method {
	case http.MethodGet:
		
	case http.MethodPost:
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func conversationHandlersWithID(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
	switch r.Method {
	case http.MethodGet:
		getConversationByIDHandler(w, r, db)
	case http.MethodPost:
		updateConversationNameByIDHandler(w, r, db)
	case http.MethodDelete:
		deleteConversationByIDHandler(w, r, db)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getConversationsHandler(w http.ResponseWriter, _ *http.Request, db *database.DataBase) {
	querry := repositories.New(db.Pool)
	conversations, err := querry.GetConversations(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)

}

func deleteConversationByIDHandler(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
	params := mux.Vars(r)
	strID := params["id"]

	UUID, err := uuid.Parse(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, rollback, commit, err := common.StartTransaction(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rollback()
	err = tx.DeleteConversationByID(context.Background(), UUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
	w.WriteHeader(http.StatusNoContent)
}

func getConversationByIDHandler(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
	params := mux.Vars(r)
	strID := params["id"]

	UUID, err := uuid.Parse(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	querry := repositories.New(db.Pool)
	conversation, err := querry.GetConversationByID(context.Background(), UUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversation)
}

func updateConversationNameByIDHandler(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
	params := mux.Vars(r)
	strID := params["id"]

	UUID, err := uuid.Parse(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var conversation repositories.UpdateConversationNameByIDParams
	err = json.NewDecoder(r.Body).Decode(&conversation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conversation.ID = UUID
	tx, rollback, commit, err := common.StartTransaction(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rollback()
	err = tx.UpdateConversationNameByID(context.Background(), conversation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
	w.WriteHeader(http.StatusNoContent)
}

func createConversationHandler(w http.ResponseWriter, r *http.Request, db *database.DataBase, storage *storage.S3Client) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Ошибка парсинга формы: "+err.Error(), http.StatusBadRequest)
		return
	}

	conversationName := r.FormValue("ConversationName")
	if conversationName == "" {
		http.Error(w, "ConversationName обязателен", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("FileUrl")
	if err != nil {
		http.Error(w, "Ошибка получения файла: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileID := uuid.New().String()
	fileKey := fmt.Sprintf("uploads/%s_%s", fileID, header.Filename)

	err = storage.UploadFile(file, fileKey)
	if err != nil {
		http.Error(w, "Ошибка загрузки файла в S3: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fileURL := fmt.Sprintf("%s/%s/%s", storage.Config.Endpoint, storage.Config.Bucket, fileKey)

	conversation := repositories.CreateConversationParams{
		ConversationName: conversationName,
		FileUrl:          fileURL,
	}

	tx, rollback, commit, err := common.StartTransaction(db)
	if err != nil {
		http.Error(w, "Ошибка начала транзакции: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rollback()

	err = tx.CreateConversation(context.Background(), conversation)
	if err != nil {
		http.Error(w, "Ошибка записи в БД: "+err.Error(), http.StatusInternalServerError)
		return
	}

	commit()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Разговор создан",
		"file_url": fileURL,
	})
}

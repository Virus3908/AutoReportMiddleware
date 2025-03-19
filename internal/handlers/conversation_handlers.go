package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"main/internal/repositories"
	"net/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (router *Router) conversationHandlers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		
	case http.MethodPost:
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (router *Router) conversationHandlersWithID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		router.getConversationByIDHandler(w, r)
	case http.MethodPost:
		router.updateConversationNameByIDHandler(w, r)
	case http.MethodDelete:
		router.deleteConversationByIDHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (router *Router) getConversationsHandler(w http.ResponseWriter, _ *http.Request) {
	querry := router.DB.NewQuerry()
	conversations, err := querry.GetConversations(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)

}

func (router *Router) deleteConversationByIDHandler(w http.ResponseWriter, r *http.Request) {
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
	err = tx.DeleteConversationByID(context.Background(), UUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
	w.WriteHeader(http.StatusNoContent)
}

func (router *Router) getConversationByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	strID := params["id"]

	UUID, err := uuid.Parse(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	querry := router.DB.NewQuerry()
	conversation, err := querry.GetConversationByID(context.Background(), UUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversation)
}

func (router *Router) updateConversationNameByIDHandler(w http.ResponseWriter, r *http.Request) {
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
	tx, rollback, commit, err := router.DB.StartTransaction()
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

func (router *Router) createConversationHandler(w http.ResponseWriter, r *http.Request) {
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

	err = router.Storage.UploadFile(file, fileKey)
	if err != nil {
		http.Error(w, "Ошибка загрузки файла в S3: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fileURL := fmt.Sprintf("%s/%s/%s", router.Storage.GetStorageEndpoint(), router.Storage.GetStorageBucket(), fileKey)

	conversation := repositories.CreateConversationParams{
		ConversationName: conversationName,
		FileUrl:          fileURL,
	}

	tx, rollback, commit, err := router.DB.StartTransaction()
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

package handlers

import (
	"encoding/json"
	"fmt"
	"main/internal/repositories"
	"net/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (router *RouterStruct) getConversationsHandler(w http.ResponseWriter, r *http.Request) {
	querry := router.DB.NewQuerry()
	conversations, err := querry.GetConversations(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)

}

func (router *RouterStruct) deleteConversationByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	strID, ok := params["id"]
	if !ok {
		http.Error(w, "missing id in request", http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(strID)
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
	err = tx.DeleteConversationByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
	w.WriteHeader(http.StatusNoContent)
}

func (router *RouterStruct) getConversationByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	strID, ok := params["id"]
	if !ok {
		http.Error(w, "missing id in request", http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	querry := router.DB.NewQuerry()
	conversation, err := querry.GetConversationByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversation)
}

func (router *RouterStruct) updateConversationNameByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	strID, ok := params["id"]
	if !ok {
		http.Error(w, "missing id in request", http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(strID)
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

	conversation.ID = id
	tx, rollback, commit, err := router.DB.StartTransaction()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rollback()
	err = tx.UpdateConversationNameByID(r.Context(), conversation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
	w.WriteHeader(http.StatusNoContent)
}

func (router *RouterStruct) createConversationHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error form parsing: "+err.Error(), http.StatusBadRequest)
		return
	}

	conversationName := r.FormValue("conversation_name")
	if conversationName == "" {
		http.Error(w, "conversation_name required", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error file recive: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileID := uuid.New().String()
	fileKey := fmt.Sprintf("uploads/%s_%s", fileID, header.Filename)

	err = router.Storage.UploadFile(file, fileKey)
	if err != nil {
		http.Error(w, "Upload to S3 error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fileURL := fmt.Sprintf("%s/%s/%s", router.Storage.GetStorageEndpoint(), router.Storage.GetStorageBucket(), fileKey)

	conversation := repositories.CreateConversationParams{
		ConversationName: conversationName,
		FileUrl:          fileURL,
	}

	tx, rollback, commit, err := router.DB.StartTransaction()
	if err != nil {
		http.Error(w, "Transaction start error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rollback()

	err = tx.CreateConversation(r.Context(), conversation)
	if err != nil {
		http.Error(w, "Error writing to db: "+err.Error(), http.StatusInternalServerError)
		return
	}

	commit()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Conversation created",
		"file_url": fileURL,
	})
}

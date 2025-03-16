package handlers

import (
	"context"
	"encoding/json"
	"main/internal/database"
	"main/internal/database/queries"
	"net/http"
	"main/internal/common"
	// "github.com/jackc/pgx/v5/pgtype"
	// "github.com/google/uuid"
	"github.com/gorilla/mux"
)

func conversationHandlers(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
	switch r.Method {
	case http.MethodGet:
		getConversationsHandler(w, r, db)
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
	querry := queries.New(db.Pool)
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

	pgUUID, err := common.StrToPGUUID(strID)
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
	err = tx.DeleteConversationByID(context.Background(), pgUUID)
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

	pgUUID, err := common.StrToPGUUID(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	querry := queries.New(db.Pool)
	conversation, err := querry.GetConversationByID(context.Background(), pgUUID)
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

	pgUUID, err := common.StrToPGUUID(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var conversation queries.UpdateConversationNameByIDParams
	err = json.NewDecoder(r.Body).Decode(&conversation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conversation.ID = pgUUID
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
package handlers

import (
	"context"
	"encoding/json"
	"main/internal/database"
	"main/internal/database/queries"
	"net/http"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func conversationHandlers(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
	switch r.Method {
	case http.MethodGet:
		getConversationsHandler(w, r, db)
	case http.MethodPost:
		createConversationsHandler(w, r, db)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func conversationHandlersWithID(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
	switch r.Method {
	case http.MethodGet:
		getConversationByIDHandler(w, r, db)
	case http.MethodPost:
		updateConversationByIDHandler(w, r, db)
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

func createConversationsHandler(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
	var conversation queries.CreateConversationParams
	err := json.NewDecoder(r.Body).Decode(&conversation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	queries := queries.New(db.Pool)
	err = queries.CreateConversation(context.Background(), conversation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func getConversationByIDHandler(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pgUUID := pgtype.UUID{Bytes: id, Valid: true}
	querry := queries.New(db.Pool)
	conversation, err := querry.GetConversationByID(context.Background(), pgUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversation)
}

func updateConversationByIDHandler(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pgUUID := pgtype.UUID{Bytes: id, Valid: true}

	var conversation queries.UpdateConversationParams
	err = json.NewDecoder(r.Body).Decode(&conversation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	conversation.ID = pgUUID
	querry := queries.New(db.Pool)
	err = querry.UpdateConversation(context.Background(), conversation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
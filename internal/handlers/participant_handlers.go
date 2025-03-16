package handlers

import (
	"context"
	"encoding/json"
	"main/internal/common"
	"main/internal/database"
	"main/internal/database/queries"
	"net/http"
	"github.com/gorilla/mux"
)

func participantHandlers(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
	switch r.Method {
	case http.MethodGet:
		getParticipantsHandler(w, r, db)
	case http.MethodPost:
		createParticipantsHandler(w, r, db)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func participantHandlersWithID(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
	switch r.Method {
	case http.MethodGet:
		getParticipantByIDHandler(w, r, db)
	case http.MethodPost:
		updateParticipantByIDHandler(w, r, db)
	case http.MethodDelete:
		deleteParticipantByIDHandler(w, r, db)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getParticipantsHandler(w http.ResponseWriter, _ *http.Request, db *database.DataBase) {
	querry := queries.New(db.Pool)
	users, err := querry.GetParticipants(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func createParticipantsHandler(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
	var user queries.CreateParticipantParams
	err := json.NewDecoder(r.Body).Decode(&user)
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
	err = tx.CreateParticipant(context.Background(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
}

func getParticipantByIDHandler(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
	params := mux.Vars(r)
	strID := params["id"]

	pgUUID, err := common.StrToPGUUID(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	querry := queries.New(db.Pool)
	user, err := querry.GetParticipantByID(context.Background(), pgUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func updateParticipantByIDHandler(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
	params := mux.Vars(r)
	strID := params["id"]

	pgUUID, err := common.StrToPGUUID(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user queries.UpdateParticipantByIDParams
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = pgUUID
	tx, rollback, commit, err := common.StartTransaction(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rollback()
	err = tx.UpdateParticipantByID(context.Background(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
}

func deleteParticipantByIDHandler(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
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
	err = tx.DeleteParticipantByID(context.Background(), pgUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
}
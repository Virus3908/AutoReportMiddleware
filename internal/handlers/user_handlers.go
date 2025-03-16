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

func userHandlers(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
	switch r.Method {
	case http.MethodGet:
		getUsersHandler(w, r, db)
	case http.MethodPost:
		createUsersHandler(w, r, db)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func userHandlersWithID(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
	switch r.Method {
	case http.MethodGet:
		getUserByIDHandler(w, r, db)
	case http.MethodPost:
		updateUserByIDHandler(w, r, db)
	case http.MethodDelete:
		deleteUserByIDHandler(w, r, db)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUsersHandler(w http.ResponseWriter, _ *http.Request, db *database.DataBase) {
	querry := queries.New(db.Pool)
	users, err := querry.GetUsers(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func createUsersHandler(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
	var user queries.CreateUserParams
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
	err = tx.CreateUser(context.Background(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
}

func getUserByIDHandler(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
	params := mux.Vars(r)
	strID := params["id"]

	pgUUID, err := common.StrToPGUUID(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	querry := queries.New(db.Pool)
	user, err := querry.GetUserByID(context.Background(), pgUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func updateUserByIDHandler(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
	params := mux.Vars(r)
	strID := params["id"]

	pgUUID, err := common.StrToPGUUID(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user queries.UpdateUserByIDParams
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
	err = tx.UpdateUserByID(context.Background(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
}

func deleteUserByIDHandler(w http.ResponseWriter, r *http.Request, db *database.DataBase) {
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
	err = tx.DeleteUserByID(context.Background(), pgUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
}
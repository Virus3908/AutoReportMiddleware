package handlers

import (
	"context"
	"encoding/json"
	"main/internal/repositories"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (router *RouterStruct) getParticipantsHandler(w http.ResponseWriter, _ *http.Request) {
	querry := router.DB.NewQuerry()
	users, err := querry.GetParticipants(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (router *RouterStruct) createParticipantHandler(w http.ResponseWriter, r *http.Request) {
	var user repositories.CreateParticipantParams
	err := json.NewDecoder(r.Body).Decode(&user)
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
	err = tx.CreateParticipant(context.Background(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
}

func (router *RouterStruct) getParticipantByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	strID := params["id"]

	UUID, err := uuid.Parse(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	querry := router.DB.NewQuerry()
	user, err := querry.GetParticipantByID(context.Background(), UUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (router *RouterStruct) updateParticipantByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	strID := params["id"]

	UUID, err := uuid.Parse(strID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user repositories.UpdateParticipantByIDParams
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = UUID
	tx, rollback, commit, err := router.DB.StartTransaction()
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

func (router *RouterStruct) deleteParticipantByIDHandler(w http.ResponseWriter, r *http.Request) {
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
	err = tx.DeleteParticipantByID(context.Background(), UUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
}
package handlers

import (
	"encoding/json"
	"main/internal/repositories"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (router *RouterStruct) getParticipantsHandler(w http.ResponseWriter, r *http.Request) {
	querry := router.DB.NewQuerry()
	users, err := querry.GetParticipants(r.Context())
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
	err = tx.CreateParticipant(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
}

func (router *RouterStruct) getParticipantByIDHandler(w http.ResponseWriter, r *http.Request) {
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
	user, err := querry.GetParticipantByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (router *RouterStruct) updateParticipantByIDHandler(w http.ResponseWriter, r *http.Request) {
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

	var user repositories.UpdateParticipantByIDParams
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = id
	tx, rollback, commit, err := router.DB.StartTransaction()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rollback()
	err = tx.UpdateParticipantByID(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
}

func (router *RouterStruct) deleteParticipantByIDHandler(w http.ResponseWriter, r *http.Request) {
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
	err = tx.DeleteParticipantByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
}
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"main/internal/repositories"
)

func (router *RouterStruct) getPromtsHandler(w http.ResponseWriter, r *http.Request) {
	querry := router.DB.NewQuerry()
	users, err := querry.GetPromts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (router *RouterStruct) getPromtByIDHandler(w http.ResponseWriter, r *http.Request) {
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
	promt, err := querry.GetPromtByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(promt)
}

func (router *RouterStruct) createPromtHandler(w http.ResponseWriter, r *http.Request) {
	var promt string
	err := json.NewDecoder(r.Body).Decode(&promt)
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
	err = tx.CreatePromt(r.Context(), promt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
}

func (router *RouterStruct) updatePromtByIDHandler(w http.ResponseWriter, r *http.Request) {
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

	var promt repositories.UpdatePromtByIDParams
	err = json.NewDecoder(r.Body).Decode(&promt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	promt.ID = id
	tx, rollback, commit, err := router.DB.StartTransaction()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rollback()
	err = tx.UpdatePromtByID(r.Context(), promt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
}

func (router *RouterStruct) deletePromtByIDHandler(w http.ResponseWriter, r *http.Request) {
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
	err = tx.DeletePromtByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commit()
}

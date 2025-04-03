package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// идея интересная, не знаю даже что я про неё думаю

func wrapperGetHandler[T any](getFn func(ctx context.Context) ([]T, error)) http.HandlerFunc { // идею понял, но нейминги странные
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := getFn(r.Context())
		if err != nil {
			respondWithError(w, err.Error(), err, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func wrapperGetByIDHandler[T any](getByIDFn func(ctx context.Context, id uuid.UUID) (T, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		strID := params["id"]
		id, err := uuid.Parse(strID)
		if err != nil {
			respondWithError(w, err.Error(), err, http.StatusBadRequest)
			return
		}

		data, err := getByIDFn(r.Context(), id)
		if err != nil {
			respondWithError(w, err.Error(), err, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func wrapperCreateHandler[T any](createFn func(ctx context.Context, payload T) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload T

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			respondWithError(w, "Invalid request body: "+err.Error(), err, http.StatusBadRequest)
			return
		}

		if err := createFn(r.Context(), payload); err != nil {
			respondWithError(w, "Failed to create: "+err.Error(), err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func wrapperUpdateHandler[T any](updateFn func(ctx context.Context, id uuid.UUID, payload T) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		strID := params["id"]
		id, err := uuid.Parse(strID)
		if err != nil {
			respondWithError(w, err.Error(), err, http.StatusBadRequest)
			return
		}

		var payload T

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			respondWithError(w, "Invalid request body: "+err.Error(), err, http.StatusBadRequest)
			return
		}

		if err := updateFn(r.Context(), id, payload); err != nil {
			respondWithError(w, "Failed to update: "+err.Error(), err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func wrapperDeleteHandler(deleteFn func(ctx context.Context, id uuid.UUID) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		strID := params["id"]
		id, err := uuid.Parse(strID)
		if err != nil {
			respondWithError(w, err.Error(), err, http.StatusBadRequest)
			return
		}

		if err := deleteFn(r.Context(), id); err != nil {
			respondWithError(w, "delete error: "+err.Error(), err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

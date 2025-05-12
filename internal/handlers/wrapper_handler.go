package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func wrapperReturningData[T any](fn func(ctx context.Context) ([]T, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := fn(r.Context())
		if err != nil {
			respondWithError(w, r, err, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func wrapperWithIDReturningData[T any](fn func(ctx context.Context, id uuid.UUID) (T, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		strID := params["id"]
		id, err := uuid.Parse(strID)
		if err != nil {
			respondWithError(w, r, err, http.StatusBadRequest)
			return
		}

		data, err := fn(r.Context(), id)
		if err != nil {
			respondWithError(w, r, err, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func wrapperWithPayload[T any](fn func(ctx context.Context, payload T) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload T

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			respondWithError(w, r, err, http.StatusBadRequest)
			return
		}

		if err := fn(r.Context(), payload); err != nil {
			respondWithError(w, r, err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func wrapperWithID(fn func(ctx context.Context, id uuid.UUID) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		strID := params["id"]
		id, err := uuid.Parse(strID)
		if err != nil {
			respondWithError(w, r, err, http.StatusBadRequest)
			return
		}

		if err := fn(r.Context(), id); err != nil {
			respondWithError(w, r, err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func wrapperWithIDAndPayload[T any](fn func(ctx context.Context, id uuid.UUID, payload T) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		strID := params["id"]
		id, err := uuid.Parse(strID)
		if err != nil {
			respondWithError(w, r, err, http.StatusBadRequest)
			return
		}

		var payload T

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			respondWithError(w, r, err, http.StatusBadRequest)
			return
		}

		if err := fn(r.Context(), id, payload); err != nil {
			respondWithError(w, r, err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
package handlers

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"
)

var ready int32

type ResponseInfo struct {
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
}

func LivenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Readiness: сервис готов только после успешного старта
func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&ready) == 1 {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Service not ready", http.StatusServiceUnavailable)
	}
}

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	info := ResponseInfo{
		Version:   "1.0.0",
		Timestamp: time.Now().Format(time.RFC3339),
		Status:    "running",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func SetReady() {
	atomic.StoreInt32(&ready, 1)
}
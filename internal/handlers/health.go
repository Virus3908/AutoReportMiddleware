package handlers

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"
)

type ResponseInfo struct {
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
}

func (router *RouterStruct) livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Readiness: сервис готов только после успешного старта
func (router *RouterStruct) readinessHandler(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&router.ready) == 1 {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Service not ready", http.StatusServiceUnavailable)
	}
}

func (router *RouterStruct) infoHandler(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&router.ready) == 1 {
		info := ResponseInfo{
			Version:   "1.0.0",
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    "running",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
	} else {
		http.Error(w, "Service not ready", http.StatusServiceUnavailable)
	}

}

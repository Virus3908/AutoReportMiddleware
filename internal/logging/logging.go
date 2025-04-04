package logging

import (
	"log"
	"net/http"
)

func LoggingMidleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		// log.Printf("From: %s", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
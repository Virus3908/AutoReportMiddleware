package logging

import (
	"log"
	"io"
	"bytes"
	"net/http"
)

func LoggingMidleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "Can't read body", http.StatusInternalServerError)
			return
		}

		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		log.Printf("Body: %s", string(bodyBytes))

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		next.ServeHTTP(w, r)
	})
}
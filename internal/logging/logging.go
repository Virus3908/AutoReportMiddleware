package logging

import (
	"log"
	"net/http"
)

func LoggingMidleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Запрос: ", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
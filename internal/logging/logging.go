package logging

import (
	"log"
	"io"
	"bytes"
	"net/http"
)

func LoggingMidleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body) // мне так нравится читать и выбрасывать нахуй сразу тело)), плюс оно не всегда в целом может быть
		// особенно мне нравится читать бесконечное огромное тело чтобы повесить все приложение либо по таймауту либо по потреблению оперативки
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "Can't read body", http.StatusInternalServerError)
			return
		}

		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		log.Printf("From: %s", r.RemoteAddr)
		// log.Printf("Body: %s", string(bodyBytes)) // видимо тело всегда есть все таки)))))))))

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		next.ServeHTTP(w, r)
	})
}
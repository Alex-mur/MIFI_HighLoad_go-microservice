package middleware

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Обертка для записи статуса ответа
		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		// Логируем запрос
		log.Info().
			Str("method", r.Method).
			Str("url", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Int("status", rw.status).
			Dur("duration", duration).
			Msg("Request completed")
	})
}

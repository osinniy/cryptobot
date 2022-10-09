package api

import (
	"net/http"
	"time"
)

// Middleware for logging requests.
// Logs request method, path, status code, hostname and response time.
func (s *Server) mLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now()

		wr := wrappedResponseWriter{w, http.StatusOK}
		next.ServeHTTP(&wr, r)

		if wr.statusCode == http.StatusMovedPermanently {
			return // don't log strict slash redirects
		}

		elapsed := time.Since(timestamp)
		s.logger.Info().
			Int("status", wr.statusCode).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Dur("elapsed_ms", elapsed).
			Msg("api request")
	})
}

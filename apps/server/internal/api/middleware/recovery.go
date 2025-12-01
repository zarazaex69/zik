package middleware

import (
	"net/http"

	"github.com/zarazaex69/zik/apps/server/internal/pkg/logger"
)

// Recovery middleware recovers from panics and logs them
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error().
					Interface("panic", err).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Msg("Panic recovered")

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": {"message": "Internal server error", "type": "internal_error", "code": 500}}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

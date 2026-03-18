package web

import (
	"log/slog"
	"net/http"
)

func (h *Handler) RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				h.Logger.ErrorContext(r.Context(), "PANIC RECOVERED",
					slog.Any("error", err),
					slog.String("path", r.URL.Path),
					slog.String("remote_addr", r.RemoteAddr),
				)

				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Internal Server Error"))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

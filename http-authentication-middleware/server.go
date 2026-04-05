package main

import (
	"log/slog"
	"net/http"
)

const validToken = "secret"

func NewServer(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux, logger)

	return mux
}

func checkAuthHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Auth-Token")

		if token != validToken {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handleHealthz(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.InfoContext(r.Context(), "health check triggered")

			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(http.StatusText(http.StatusOK))); err != nil {
				logger.ErrorContext(r.Context(), "failed to write response", "error", err)
			}
		},
	)
}

func handlerSecure(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.InfoContext(r.Context(), "hit secure endpoint")

			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(http.StatusText(http.StatusOK))); err != nil {
				logger.ErrorContext(r.Context(), "failed to write response", "error", err)
			}
		},
	)
}

func addRoutes(mux *http.ServeMux, logger *slog.Logger) {
	mux.Handle("/healthz", handleHealthz(logger))
	mux.Handle("/secure", checkAuthHeader(handlerSecure(logger)))
}

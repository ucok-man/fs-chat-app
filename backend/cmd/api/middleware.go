package main

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *application) recover(next http.Handler) http.Handler {
	return middleware.Recoverer(next)
}

func (app *application) cors(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   app.config.cors.trustedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Access-Control-Request-Method"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		// Caching duration for preflight requests (in seconds)
		MaxAge: 60,
	})(next)
}

package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

const BASE_PATH = "/api/v1"

func (app *application) routes() http.Handler {
	router := chi.NewRouter()
	router.NotFound(http.HandlerFunc(app.errNotFoundResponse))
	router.MethodNotAllowed(http.HandlerFunc(app.errMethodNotAllowedResponse))

	// Basic Middleware
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   app.config.cors.trustedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Access-Control-Request-Method"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		// Caching duration for preflight requests (in seconds)
		MaxAge: 60,
	}))

	apirouter := chi.NewRouter()
	apirouter.Get("/healthcheck", app.healthcheckHandler)

	router.Mount("/api/v1", apirouter)
	return router
}

package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	apiv1 := chi.NewRouter()
	apiv1.NotFound(http.HandlerFunc(app.errNotFoundResponse))
	apiv1.MethodNotAllowed(http.HandlerFunc(app.errMethodNotAllowedResponse))

	apiv1.Get("/healthcheck", app.healthcheckHandler)

	// public routes
	apiv1.Group(func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signup", app.signup)
			r.Post("/signin", app.signin)
			r.Post("/logout", app.logout)
		})
	})

	// private routes

	router := chi.NewRouter().With(app.recover, app.cors)
	router.Mount("/api/v1", apiv1)
	return router
}

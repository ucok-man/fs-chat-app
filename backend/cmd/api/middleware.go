package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ucok-man/fs-chat-app-backend/internal/data"
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

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwttoken")
		if err != nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &JwtTokenClaim{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(app.config.jwt.secret), nil
		})

		if err != nil || !token.Valid {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		claim, ok := token.Claims.(*JwtTokenClaim)
		if !ok {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.models.User.GetById(claim.Uid)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.errServerResponse(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)
		next.ServeHTTP(w, r)

	})
}

package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/ucok-man/fs-chat-app-backend/cmd/api/dto"
	"github.com/ucok-man/fs-chat-app-backend/internal/data"
)

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	var input dto.ReqSignupDto
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errBadRequestResponse(w, r, err)
		return
	}

	if err := input.Validate(); err != nil {
		app.errFailedValidationResponse(w, r, err)
		return
	}

	user := &data.User{
		FullName: input.FullName,
		Email:    input.Email,
	}

	if err := user.Password.Set(input.Password); err != nil {
		app.errServerResponse(w, r, err)
		return
	}

	err = app.models.User.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			app.errFailedValidationResponse(w, r, map[string]string{
				"email": "a user with this email address already exists",
			})
		default:
			app.errServerResponse(w, r, err)
		}
		return
	}

	token, err := app.generateToken(user.ID)
	if err != nil {
		app.errServerResponse(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwttoken",
		HttpOnly: true, // Ensures the cookie is not accessible via JavaScript
		SameSite: http.SameSiteStrictMode,
		Secure:   app.config.env == "production",
		MaxAge:   int((7 * 24 * time.Hour).Seconds()), // 7 days
		Value:    token,
	})

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.errServerResponse(w, r, err)
		return
	}
}

func (app *application) signin(w http.ResponseWriter, r *http.Request) {
	var input dto.ReqSigninDto
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errBadRequestResponse(w, r, err)
		return
	}

	if err := input.Validate(); err != nil {
		app.errFailedValidationResponse(w, r, err)
		return
	}

	user, err := app.models.User.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.errInvalidCredentialsResponse(w, r)
		default:
			app.errServerResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.errServerResponse(w, r, err)
		return
	}

	if !match {
		app.errInvalidCredentialsResponse(w, r)
		return
	}

	token, err := app.generateToken(user.ID)
	if err != nil {
		app.errServerResponse(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwttoken",
		HttpOnly: true, // Ensures the cookie is not accessible via JavaScript
		SameSite: http.SameSiteStrictMode,
		Secure:   app.config.env == "production",
		MaxAge:   int((7 * 24 * time.Hour).Seconds()), // 7 days
		Value:    token,
	})

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.errServerResponse(w, r, err)
		return
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "jwttoken",
		MaxAge: -1,
	})

	err := app.writeJSON(w, http.StatusNoContent, envelope{}, nil)
	if err != nil {
		app.errServerResponse(w, r, err)
		return
	}
}

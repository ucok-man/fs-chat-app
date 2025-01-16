package main

import (
	"errors"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	app.errServerResponse(w, r, errors.New("AAA"))
	return
	// env := envelope{
	// 	"status": "available",
	// 	"system_info": map[string]string{
	// 		"environment": app.config.env,
	// 		"version":     version,
	// 	},
	// }

	// err := app.writeJSON(w, http.StatusOK, env, nil)
	// if err != nil {
	// 	app.errServerResponse(w, r, err)
	// 	return
	// }
}

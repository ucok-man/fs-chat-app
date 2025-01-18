package main

import (
	"errors"
	"net/http"

	"github.com/ucok-man/fs-chat-app-backend/internal/data"
	"github.com/ucok-man/fs-chat-app-backend/internal/media"
)

func (app *application) uploadProfile(w http.ResponseWriter, r *http.Request) {
	filehandler, err := app.readImageFile(r, "profile")
	if err != nil {
		app.errBadRequestResponse(w, r, err)
		return
	}

	file, err := filehandler.Open()
	if err != nil {
		app.errServerResponse(w, r, err)
		return
	}
	defer file.Close()

	user := app.contextGetUser(r)

	imgurl, err := app.media.Upload(
		file,
		media.UploadWithFolder("fs-chat-app/profile/"),
		media.UploadWithReplaceable(user.ID),
	)
	if err != nil {
		app.errServerResponse(w, r, err)
		return
	}

	user.ProfilePic = imgurl
	if err := app.models.User.Update(user); err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.errEditConflictResponse(w, r)
		default:
			app.errServerResponse(w, r, err)

		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.errServerResponse(w, r, err)
		return
	}
}

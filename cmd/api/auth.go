package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/atomicmeganerd/gopher-social/internal/mailer"
	"github.com/atomicmeganerd/gopher-social/internal/store"
	"github.com/google/uuid"
)

type RegisteredUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a new user
//	@Description	Registers a new user with the provided details
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisteredUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken			"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/users/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var payload RegisteredUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Hash the password before storing it

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	plainToken := uuid.New().String()
	// Hash the token to store in the database (this encrypts it)
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	if err := app.store.Users.CreateAndInvite(
		r.Context(), user, hashToken, app.config.mail.exp,
	); err != nil {
		switch err {
		case store.ErrDuplicateUsername:
			app.badRequestError(w, r, err)
		case store.ErrDuplicateEmail:
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)
	isProdEnv := app.config.env == "production"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	status, err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)

	if err != nil {
		if err := app.store.Users.Delete(r.Context(), user.ID); err != nil {
			app.logger.Error("error deleting user", "user", user)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	app.logger.Info("Email sent", "status code", status)

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
	}
}

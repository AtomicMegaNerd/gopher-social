package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/atomicmeganerd/gopher-social/internal/mailer"
	"github.com/atomicmeganerd/gopher-social/internal/store"
	"github.com/golang-jwt/jwt/v5"
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
		Role: store.Role{
			Name: "user",
		},
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	plainToken := uuid.New().String()
	// Hash the token to store in the database (this encrypts it)
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	if err := app.dbStore.Users.CreateAndInvite(
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

	status, err := app.mailer.Send(
		mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv,
	)

	if err != nil {
		if err := app.dbStore.Users.Delete(r.Context(), user.ID); err != nil {
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

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// createTokenHandler godoc
//
//	@Summary		Creates a token
//	@Description	Creates a token for a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"User credentials"
//	@Success		200		{string}	string					"Token"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserTokenPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Check that jwt is configured
	if app.config.auth.jwtToken.tokenHost == "" || app.config.auth.jwtToken.secret == "" {
		app.internalServerError(w, r, errors.New("jwt is not configured on this app instance"))
		return
	}

	user, err := app.dbStore.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			// NOTE: do not use 404 here due to enumeration attacks!
			app.unauthorizedError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

	}

	if err := user.Password.Compare(payload.Password); err != nil {
		app.unauthorizedError(w, r, errors.New("invalid password"))
	}

	// NOTE: Check the docs on JWT to see what claims can be setup
	// See: https://auth0.com/docs/secure/tokens/json-web-tokens/json-web-token-claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.jwtToken.expiry).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.jwtToken.tokenHost,
		"aud": app.config.auth.jwtToken.tokenHost,
	}
	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// send the token to the client
	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)
	}
}

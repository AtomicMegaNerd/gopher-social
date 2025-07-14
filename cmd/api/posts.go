package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/atomicmeganerd/rcd-gopher-social/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		slog.Error("failed to read JSON", "error", err)
		writeJSONError(w, http.StatusBadRequest, err.Error())
	}
	userId := 1 // TODO: This should be replaced with actual user ID extraction logic
	ctx := r.Context()

	post := &store.Post{
		UserID:  int64(userId),
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
	}

	if err := app.store.Posts.Create(ctx, post); err != nil {
		slog.Error("failed to create post", "error", err)
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("failed to create post %s", err.Error()))
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		slog.Error("failed to write JSON response", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {

	postIDRaw := chi.URLParam(r, "postID")
	if postIDRaw == "" {
		slog.Error("Missing argument postID")
		writeJSONError(w, http.StatusBadRequest, "you did not include a postID")
		return
	}

	postID, err := strconv.ParseInt(postIDRaw, 10, 64)
	if err != nil {
		slog.Error("invalid request postID must be integer", "postID", postIDRaw)
		writeJSONError(w, http.StatusBadRequest, "invalid postID must be integer")
		return
	}
	ctx := r.Context()

	post, err := app.store.Posts.GetByID(ctx, postID)
	if err != nil {
		slog.Error("unable to load post", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Unable to load post")
		return
	}

	_ = writeJSON(w, http.StatusOK, post)
}

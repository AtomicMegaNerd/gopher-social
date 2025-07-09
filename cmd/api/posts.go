package main

import (
	"net/http"

	"github.com/atomicmeganerd/rcd-gopher-social/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
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
		writeJSONError(w, http.StatusInternalServerError, "failed to create post")
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

package main

import "net/http"

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Add pagination, search, filters. etc.
	feed, err := app.store.Posts.GetUserFeed(r.Context(), int64(41)) // Temporary user ID
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}

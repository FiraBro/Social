package main

import (
	"net/http"
)

type createPostPayload struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var postPayload createPostPayload

	// Step 1: Decode JSON from the request body into postPayload
	if err := readJSON(w, r, &postPayload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Step 2: Create a new Post struct (from your store or model)
	post := &store.Post{
		Title:   postPayload.Title,
		Content: postPayload.Content,
		UserID:  1, // temporarily hardcoded for testing
	}

	// Step 3: Insert the post into the database
	ctx := r.Context()
	if err := app.store.Post.Create(ctx, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Step 4: Send back the created post as a response
	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

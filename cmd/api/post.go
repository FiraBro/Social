package main

import (
	"log"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/FiraBro/social/internal/store"
)

const temporaryPostUserID int64 = 1

type createPostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var postPayload createPostPayload
	if err := readJSON(w, r, &postPayload); err != nil {
		_ = writeJSONError(w, http.StatusBadRequest, "invalid JSON request body")
		return
	}

	if message := validateCreatePostPayload(postPayload); message != "" {
		_ = writeJSONError(w, http.StatusBadRequest, message)
		return
	}

	if postPayload.Tags == nil {
		postPayload.Tags = []string{}
	}
	for i := range postPayload.Tags {
		postPayload.Tags[i] = strings.TrimSpace(postPayload.Tags[i])
	}

	post := &store.Post{
		Title:   strings.TrimSpace(postPayload.Title),
		Content: postPayload.Content,
		UserID:  temporaryPostUserID,
		Tags:    postPayload.Tags,
	}

	if err := app.store.Post.Create(r.Context(), post); err != nil {
		log.Printf("create post: %v", err)
		_ = writeJSONError(w, http.StatusInternalServerError, "could not create post")
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		log.Printf("write create post response: %v", err)
	}
}

func validateCreatePostPayload(payload createPostPayload) string {
	titleLength := utf8.RuneCountInString(strings.TrimSpace(payload.Title))
	if titleLength == 0 {
		return "title is required"
	}
	if titleLength > 255 {
		return "title must be 255 characters or fewer"
	}

	contentLength := utf8.RuneCountInString(strings.TrimSpace(payload.Content))
	if contentLength == 0 {
		return "content is required"
	}
	if contentLength > 10000 {
		return "content must be 10000 characters or fewer"
	}

	if len(payload.Tags) > 20 {
		return "tags must contain 20 items or fewer"
	}
	for _, tag := range payload.Tags {
		tagLength := utf8.RuneCountInString(strings.TrimSpace(tag))
		if tagLength == 0 {
			return "tags must not contain empty values"
		}
		if tagLength > 50 {
			return "tags must be 50 characters or fewer"
		}
	}

	return ""
}

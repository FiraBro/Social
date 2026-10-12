package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/FiraBro/social/internal/store"
)

type fakePostRepository struct {
	called bool
	post   *store.Post
	err    error
	ctxErr error
}

func (f *fakePostRepository) Create(ctx context.Context, post *store.Post) error {
	f.called = true
	f.ctxErr = ctx.Err()
	if f.ctxErr != nil {
		return f.ctxErr
	}
	if f.err != nil {
		return f.err
	}

	post.ID = "generated-id"
	post.Status = "draft"
	post.CreatedAt = time.Date(2026, time.October, 12, 12, 0, 0, 0, time.UTC)
	post.UpdatedAt = post.CreatedAt
	f.post = post
	return nil
}

func TestCreatePostHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		repoErr    error
		wantStatus int
		wantCall   bool
	}{
		{
			name:       "creates post",
			body:       `{"title":" First post ","content":"Learning Go","tags":["golang","backend"]}`,
			wantStatus: http.StatusCreated,
			wantCall:   true,
		},
		{
			name:       "malformed JSON",
			body:       `{"title":`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "trailing JSON value",
			body:       `{"title":"ok","content":"ok"} {}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "unknown field",
			body:       `{"title":"title","content":"content","user_id":9}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing title",
			body:       `{"content":"content"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "blank content",
			body:       `{"title":"title","content":"  "}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "title too long",
			body:       `{"title":"` + strings.Repeat("a", 256) + `","content":"content"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "too many tags",
			body:       `{"title":"title","content":"content","tags":[` + strings.Repeat(`"tag",`, 20) + `"tag"]}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "repository error is hidden",
			body:       `{"title":"title","content":"content"}`,
			repoErr:    errors.New("database password leaked"),
			wantStatus: http.StatusInternalServerError,
			wantCall:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &fakePostRepository{err: tt.repoErr}
			app := &application{store: store.Storage{Post: repository}}
			request := httptest.NewRequest(http.MethodPost, "/v1/posts", strings.NewReader(tt.body))
			response := httptest.NewRecorder()

			app.mount().ServeHTTP(response, request)

			if response.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body: %s", response.Code, tt.wantStatus, response.Body.String())
			}
			if repository.called != tt.wantCall {
				t.Fatalf("repository called = %t, want %t", repository.called, tt.wantCall)
			}

			if tt.name == "creates post" {
				if repository.post.Title != "First post" {
					t.Errorf("stored title = %q, want trimmed title", repository.post.Title)
				}
				if repository.post.UserID != temporaryPostUserID {
					t.Errorf("user ID = %d, want temporary ID %d", repository.post.UserID, temporaryPostUserID)
				}
				if len(repository.post.Tags) != 2 || repository.post.Tags[0] != "golang" {
					t.Errorf("stored tags = %v, want [golang backend]", repository.post.Tags)
				}
				if !strings.Contains(response.Body.String(), `"id":"generated-id"`) {
					t.Errorf("response does not contain generated post ID: %s", response.Body.String())
				}
			}
			if tt.name == "repository error is hidden" && strings.Contains(response.Body.String(), "database password leaked") {
				t.Errorf("response exposed repository error: %s", response.Body.String())
			}
		})
	}
}

func TestCreatePostHandlerPassesCanceledContext(t *testing.T) {
	repository := &fakePostRepository{}
	app := &application{store: store.Storage{Post: repository}}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	request := httptest.NewRequest(http.MethodPost, "/v1/posts", strings.NewReader(`{"title":"title","content":"content"}`)).WithContext(ctx)
	response := httptest.NewRecorder()

	app.createPostHandler(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
	if !errors.Is(repository.ctxErr, context.Canceled) {
		t.Fatalf("repository context error = %v, want context.Canceled", repository.ctxErr)
	}
}

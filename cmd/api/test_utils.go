package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/atomicmeganerd/gopher-social/internal/auth"
	"github.com/atomicmeganerd/gopher-social/internal/ratelimiter"
	"github.com/atomicmeganerd/gopher-social/internal/store"
	"github.com/atomicmeganerd/gopher-social/internal/store/cache"
	"github.com/lmittmann/tint"
)

func newTestApp(t *testing.T, cfg config) *application {
	t.Helper()

	handler := tint.NewHandler(os.Stderr, &tint.Options{
		Level:     slog.LevelInfo,
		AddSource: true,
	})
	logger := slog.New(handler)

	mockStore := store.NewMockStore()
	mockCache := cache.NewMockStore()
	mockAuth := &auth.TestAuthenticator{}

	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	return &application{
		logger:        logger,
		dbStore:       mockStore,
		cacheStore:    mockCache,
		authenticator: mockAuth,
		config:        cfg,
		rateLimiter:   rateLimiter,
	}
}

func execMockRequests(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected response code %d but got %d", expected, actual)
	}
}

package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/atomicmeganerd/gopher-social/internal/store"
	"github.com/atomicmeganerd/gopher-social/internal/store/cache"
	"github.com/lmittmann/tint"
)

func newTestApp(t *testing.T) *application {
	t.Helper()

	handler := tint.NewHandler(os.Stderr, &tint.Options{
		Level:     slog.LevelInfo,
		AddSource: true,
	})
	logger := slog.New(handler)

	mockStore := store.NewMockStore()
	mockCache := cache.NewMockStore()

	return &application{
		logger:     logger,
		dbStore:    mockStore,
		cacheStore: mockCache,
	}
}

func execMockRequests(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

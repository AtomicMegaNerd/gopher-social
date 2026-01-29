package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/atomicmeganerd/gopher-social/internal/store/cache"
	"github.com/stretchr/testify/mock"
)

func TestGetUser(t *testing.T) {
	withRedis := config{
		cache: cacheConfig{
			enabled: true,
		},
	}

	app := newTestApp(t, withRedis)
	mux := app.mount()

	testToken, err := app.authenticator.GenerateToken(nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := execMockRequests(req, mux)
		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should allow authenticated requests", func(t *testing.T) {

		mockCacheStore := app.cacheStore.Users.(*cache.MockUsersCacheStorage)
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		mockCacheStore.On("Get", int64(1)).Return(nil, nil).Twice()
		mockCacheStore.On("Set", mock.Anything).Return(nil)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testToken))
		rr := execMockRequests(req, mux)
		checkResponseCode(t, http.StatusOK, rr.Code)

		mockCacheStore.Calls = nil // Reset mock expectations
	})

	t.Run(
		"should hit the cache first and if not exists it sets the user in the cache",
		func(t *testing.T) {
			mockCacheStore := app.cacheStore.Users.(*cache.MockUsersCacheStorage)

			mockCacheStore.On("Get", int64(42)).Return(nil, nil)
			mockCacheStore.On("Get", int64(1)).Return(nil, nil)
			mockCacheStore.On("Set", mock.Anything, mock.Anything).Return(nil)

			req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testToken))
			rr := execMockRequests(req, mux)

			checkResponseCode(t, http.StatusOK, rr.Code)

			mockCacheStore.AssertNumberOfCalls(t, "Get", 2)
			mockCacheStore.Calls = nil

		})

	t.Run("should NOT hit the cache if it is not enabled", func(t *testing.T) {

		withRedis := config{
			cache: cacheConfig{
				enabled: false,
			},
		}

		app := newTestApp(t, withRedis)
		mux := app.mount()

		mockCacheStore := app.cacheStore.Users.(*cache.MockUsersCacheStorage)

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)
		rr := execMockRequests(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)

		mockCacheStore.AssertNotCalled(t, "Get")
		mockCacheStore.Calls = nil
	})
}

package main

import (
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	app := newTestApp(t)
	mux := app.mount()

	t.Run("should not allow unauthenticated requests", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, "v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := execMockRequests(req, mux)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected response code %d but got %d", http.StatusUnauthorized, rr.Code)
		}
	})
}

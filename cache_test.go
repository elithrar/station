package station

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := Cache()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	handler.ServeHTTP(rr, r)

	if rr.Header().Get(cacheControl) != fmt.Sprintf("%s%d", cacheControlValue,
		week) {
		t.Fatalf("Cache-Control header not set correctly.")
	}

	if rr.Header().Get(expires) != "" {
		t.Fatalf("Expires header incorrectly set (not empty for %v request)",
			r.Proto)
	}

	if rr.Header().Get(pragma) != "" {
		t.Fatalf("Pragma header incorrectly set (not empty)")
	}
}

func TestCacheExpires(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	r.ProtoMinor = 0
	handler := Cache()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	handler.ServeHTTP(rr, r)

	if rr.Header().Get(expires) != time.Now().Add(time.Duration(week)).Format(time.RFC1123) {
		t.Fatalf("Expires header invalid: got %v", rr.Header().Get(expires))
	}
}

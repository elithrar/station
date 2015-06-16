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

	var hour = time.Hour * 1

	opts := CacheOptions{
		MaxAge: time.Hour * 1,
	}

	rr := httptest.NewRecorder()

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := Cache(opts)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	handler.ServeHTTP(rr, r)

	if rr.Header().Get(cacheControl) != fmt.Sprintf("%s%.0f", cacheControlValue,
		hour.Seconds()) {
		t.Fatalf("Cache-Control header not set correctly.")
	}

	if rr.Header().Get(pragma) != "" {
		t.Fatalf("Pragma header incorrectly set (not empty)")
	}
}

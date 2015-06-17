package station

import (
	"fmt"
	"net/http"
	"time"
)

const (
	cacheControl      = "Cache-Control"
	cacheControlValue = "public, must-revalidate, max-age="
	expires           = "Expires"
	pragma            = "Pragma"
)

var month = int64(86400 * 30)

type cache struct {
	h    http.Handler
	opts cacheOptions
}

// CacheOptions stores configuration options for cache headers.
type cacheOptions struct {
	maxAge         int64
	expires        time.Time
	mustRevalidate bool
}

// Satisfies the http.Handler interface for cache. Behaviour is defined as per
// RFC2616 and https://www.mnot.net/cache_docs/#CACHE-CONTROL
func (c cache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set Cache-Control header (HTTP/1.1 clients)
	w.Header().Set(cacheControl,
		fmt.Sprintf("%s%d", cacheControlValue, c.opts.maxAge))

	// Set the Expires header (HTTP/1.0 clients only)
	if r.ProtoMinor == 0 {
		w.Header().Set(expires, time.Now().Add(time.Duration(
			c.opts.maxAge)).Format(time.RFC1123))
	}

	// Unset any Pragma header.
	w.Header().Del(pragma)

	// Call the next handler
	c.h.ServeHTTP(w, r)
}

// MaxAge sets the duration to cache objects for in seconds.
func MaxAge(age int64) func(*cache) {
	return func(c *cache) {
		c.opts.maxAge = age
	}
}

// MustRevalidate sets the 'must-revalidate' Cache-Control parameter to request
// that the client check the If-Modified-Since, Last-Modified or ETag headers
// before serving a cached file to the user. Defaults to false, but is recommended.
func MustRevalidate() func(*cache) {
	return func(c *cache) {
		c.opts.mustRevalidate = true
	}
}

func parseCache(h http.Handler, options ...func(*cache)) *cache {
	c := &cache{h: h}

	for _, option := range options {
		option(c)
	}

	return c
}

// Cache provides HTTP middleware for setting client-side caching headers for
// HTTP resources. These headers are commonly used to set far-future dates for
// static assets to minimise additional HTTP requests on repeat visits.
func Cache(options ...func(*cache)) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		c := parseCache(h, options...)

		// Default cache duration is one month
		if c.opts.maxAge == 0 {
			c.opts.maxAge = month
		}

		return c
	}
}

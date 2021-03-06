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

var week = int64(86400 * 7)

type cache struct {
	h    http.Handler
	opts cacheOptions
}

// CacheOption represents an option for configuring the Cache handler.
type CacheOption func(*cache)

// CacheOptions stores configuration options for cache headers.
type cacheOptions struct {
	maxAge int64
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
func MaxAge(age int64) CacheOption {
	return func(c *cache) {
		c.opts.maxAge = age
	}
}

func parseCache(h http.Handler, options ...CacheOption) *cache {
	c := &cache{h: h}

	for _, option := range options {
		option(c)
	}

	return c
}

// Cache provides HTTP middleware for setting client-side caching headers for
// HTTP resources. These headers are commonly used to set far-future dates for
// static assets to minimise additional HTTP requests on repeat visits.
func Cache(options ...CacheOption) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		c := parseCache(h, options...)

		// Default cache duration is one week
		if c.opts.maxAge == 0 {
			c.opts.maxAge = week
		}

		return c
	}
}

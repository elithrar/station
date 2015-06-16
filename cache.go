package station

import (
	"fmt"
	"net/http"
	"time"
)

const (
	cacheControl      = "Cache-Control"
	cacheControlValue = "public, must-revalidate, max-age="
	epoch             = time.Duration(0)
	expires           = "Expires"
	pragma            = "Pragma"
	month             = 2592000
)

type cache struct {
	h    http.Handler
	opts CacheOptions
}

// CacheOptions stores configuration options for cache headers.
type CacheOptions struct {
	// Set the duration to cache objects for
	MaxAge time.Duration
	// Set the must-revalidate parameter to request that the client check the
	// If-Modified-Since, Last-Modified or ETag headers before serving a cached
	// file. Defaults to false, but is recommended.
	MustRevalidate bool
}

// Satisfies the http.Handler interface for cache. Behaviour is defined as per
// RFC2616 and https://www.mnot.net/cache_docs/#CACHE-CONTROL
func (c cache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set Cache-Control header (HTTP/1.1 clients)
	w.Header().Set(cacheControl,
		fmt.Sprintf("%s%.0f", cacheControlValue, c.opts.MaxAge.Seconds()))

	// Set the Expires header (HTTP/1.0 clients only)
	if r.ProtoMinor == 0 {
		w.Header().Set(expires,
			time.Now().Add(c.opts.MaxAge).Format(time.RFC3339))
	}

	// Unset any Pragma header.
	w.Header().Del(pragma)

	// Call the next handler
	c.h.ServeHTTP(w, r)
}

// Cache provides HTTP middleware for setting client-side caching headers for
// HTTP resources. These headers are commonly used to set far-future dates for
// static assets to minimise additional HTTP requests on repeat visits.
func Cache(opts CacheOptions) func(http.Handler) http.Handler {
	// Default to one month if options are unset.
	if opts.MaxAge == epoch {
		opts.MaxAge = time.Duration(month)
	}

	fn := func(h http.Handler) http.Handler {
		return cache{h, opts}
	}

	return fn
}

// Package station provides HTTP static file serving & caching middleware.
package station

import (
	"net/http"
	"os"
	"path/filepath"
)

type static struct {
	dir  string
	h    http.Handler
	opts StaticOptions
}

// StaticOptions sets the options for serving static files using ServeStatic.
type StaticOptions struct {
	// Turn directory listings on (i.e. show all files in a directory).
	ListDir bool
	// NotFound is called when using ServeStatic. Defaults to
	// http.NotFoundHandler if not provided.
	NotFoundHandler http.Handler
}

// Satifies http.Handler for static. The Content-Type header is automatically
// set by http.ServeFile based on Go's content type detection.
func (s static) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fname := filepath.Join(s.dir, r.URL.Path)

	// Check if the file exists. If not, call the wrapped handler.
	f, err := os.Stat(fname)
	if err != nil {
		s.h.ServeHTTP(w, r)
		return
	}

	// Don't show directory listings if the option isn't set.
	if f.IsDir() && !s.opts.ListDir {
		s.h.ServeHTTP(w, r)
		return
	}

	// http.ServeFile sets Last-Modified headers based on modtime for us.
	http.ServeFile(w, r, fname)
}

// Static provides HTTP middleware that serves static assets from the directory
// provided. If the file doesn't exist, it calls the wrapped handler/router.
// This is useful when you want static files in a directory to be served as a
// first priority (e.g. favicon.ico, stylesheets, etc.) across an entire router.
func Static(dir string, opts StaticOptions) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return static{dir, h, opts}
	}
}

// Serve is a handler that serves static files from the directory
// provided. If the file doesn't exist, it calls opts.NotFound.
func Serve(dir string, opts StaticOptions) http.Handler {
	if opts.NotFoundHandler == nil {
		opts.NotFoundHandler = http.NotFoundHandler()
	}

	return Static(dir, opts)(opts.NotFoundHandler)
}

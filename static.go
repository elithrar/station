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
	opts staticOptions
}

// staticOptions sets the options for serving static files.
type staticOptions struct {
	// Turn directory listings on (i.e. show all files in a directory).
	dirList bool
	// NotFound is called when using ServeStatic. Defaults to
	// http.NotFoundHandler if not provided.
	notFoundHandler http.Handler
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
	if f.IsDir() && !s.opts.dirList {
		s.h.ServeHTTP(w, r)
		return
	}

	// http.ServeFile sets Last-Modified headers based on modtime for us.
	http.ServeFile(w, r, fname)
}

// DirList turns directory listings 'on' (off by default).
func DirList() func(*static) {
	return func(s *static) {
		s.opts.dirList = true
	}
}

// NotFoundHandler sets a custom http.Handler to be called when using the Serve
// handler. Set this to serve 'pretty' HTTP 404 pages or re-directs.
// elsewhere.
func NotFoundHandler(h http.Handler) func(*static) {
	return func(s *static) {
		s.opts.notFoundHandler = h
	}
}

func parseStatic(dir string, h http.Handler, options ...func(*static)) *static {
	s := &static{
		dir: dir,
		h:   h,
	}

	for _, option := range options {
		option(s)
	}

	return s
}

// Static provides HTTP middleware that serves static assets from the directory
// provided. If the file doesn't exist, it calls the wrapped handler/router.
// This is useful when you want static files in a directory to be served as a
// first priority (e.g. favicon.ico, stylesheets, etc.) across an entire router.
func Static(dir string, options ...func(*static)) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return parseStatic(dir, h, options...)
	}
}

// Serve is a handler that serves static files from the directory
// provided. If the file doesn't exist, it calls the currently configured
// NotFoundHandler (defaults to http.NotFoundHandler).
func Serve(dir string, options ...func(*static)) http.Handler {
	s := parseStatic(dir, nil, options...)

	// Use the built-in HTTP 404 Not Found handler from net/http if unset
	if s.opts.notFoundHandler == nil {
		s.opts.notFoundHandler = http.NotFoundHandler()
	}
	// We call the notFoundHandler as the wrapped handler when not operating
	// as middleware.
	s.h = s.opts.notFoundHandler

	return s
}

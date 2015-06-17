# station [![GoDoc](https://godoc.org/github.com/elithrar/station?status.svg)](https://godoc.org/github.com/elithrar/station) [![Build Status](https://travis-ci.org/elithrar/station.svg)](https://travis-ci.org/elithrar/station)

station provides simple static file serving & caching middleware for Go HTTP
applications.

* `Serve` is a handler that will serve files below the supplied path, with 
 directory listings an optional extra. 
* `Static` is HTTP middleware that serves files (if they exist) before passing
  the request to the rest of your application.
* `Cache` is HTTP middleware that sets a number of useful HTTP caching headers, 
 including `Cache-Control` and `Expires`.

See below for a full set of examples.

## Examples

*station* can be used with Go's [net/http](http://golang.org/pkg/net/http/)
package, with the [Goji micro-framework](https://github.com/zenazn/goji) or with
[gorilla/mux](http://www.gorillatoolkit.org/pkg/mux).

Here's an example of using it with plain old `net/http`:

```go
import (
    "net/http"

    "github.com/elithrar/station"
)

func main() {
    r := http.NewServeMux()
    opts := StaticOptions{
        // Directory listings are off by default. Let's turn them on.
        ListDir: true,
    }

    // IndexHandler is just a func(w http.ResponseWriter, r *http.Request) here.
    r.HandleFunc("/", IndexHandler)
    r.Handle("/static/", station.Serve("/Users/matt/Desktop/static",
    opts))

    http.ListenAndServe(":8000", r)
}
```

... or with Goji:

```go
import (
    "net/http"

    "github.com/elithrar/station"
    "github.com/zenazn/goji"
)

func main() {
    goji.Get("/", IndexHandler)
    goji.Get("/static/*", station.Serve("/Users/matt/Desktop/static",
    StaticOptions{}))
    goji.Serve()
}
```

Pretty easy, huh? You can also use *station* as HTTP middleware, where
it will check for and serve a static file before falling back to your other
handlers. Here's how to do that:

```go
import (
    "net/http"

    "github.com/elithrar/station"
)

func main() {
    r := http.NewServeMux()
    // IndexHandler is just a func(w http.ResponseWriter, r *http.Request) here.
    r.HandleFunc("/", IndexHandler)

    // We wrap our top-level router with the Static middleware. When you request
    // a path like http://example.com/static/style.css, the middleware will
    // check if that file exists before passing the request off to your router.
    // PS: If you're using Goji, just call goji.Use(station.Static(opts))
    http.ListenAndServe(":8000", station.Static("/Users/matt/Desktop/static", StaticOptions{})(r))
}
```

If you still have questions about how to use station, raise an issue with 
your code (as it stands) and the error message (if any).

PRs for bugs and new features are welcomed, with API stability the #1 priority.

# License

BSD licensed. See the LICENSE file for details.

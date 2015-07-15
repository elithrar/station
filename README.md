# station
[![GoDoc](https://godoc.org/github.com/elithrar/station?status.svg)](https://godoc.org/github.com/elithrar/station) [![Build Status](https://travis-ci.org/elithrar/station.svg?branch=master)](https://travis-ci.org/elithrar/station)

station provides simple static file serving & caching middleware for Go HTTP
applications. It makes it easy to turn off directory listings (off by default,
actually) and set HTTP caching headers on any routes or routers you need to.

* `Serve` is a handler that will serve files below the supplied path, with
 directory listings an optional extra. 
* `Static` is HTTP middleware that serves files (if they exist) before passing
  the request to the rest of your application.
* `Cache` is HTTP middleware that sets a number of useful HTTP caching headers, 
 including `Cache-Control` and `Expires`.

See below for a full set of examples.

## Examples

### Static File Serving

*station* can be used with Go's [net/http](http://golang.org/pkg/net/http/)
package, with the [Goji micro-framework](https://github.com/zenazn/goji) and with
[gorilla/mux](http://www.gorillatoolkit.org/pkg/mux).

Here's an example of using it with plain old `net/http`:

```go
import (
    "net/http"

    "github.com/elithrar/station"
)

func main() {
    r := http.NewServeMux()

    // IndexHandler is just a func(w http.ResponseWriter, r *http.Request) here.
    r.HandleFunc("/", IndexHandler)

    // Directory listings are off by default. Let's turn them on.
    r.Handle("/static/", station.Serve("/Users/matt/Desktop/static",
    station.DirList())

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

    // Just pass the path to your static files if you want to use the defaults.
    goji.Get("/static/*", station.Serve("/Users/matt/Desktop/static"))
    goji.Serve()
}
```

Pretty easy, huh? You can also use *station* as HTTP middleware, where
it will check for and serve a static file from `yourdomain.com/your/static/path` 
before falling back to your router. Here's how to do that:

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
    http.ListenAndServe(":8000", station.Static("/Users/matt/Desktop/static/")(r))
}
```

### Caching Middleware

Using `station.Cache` is also pretty easy:

```go
import (
    "net/http"

    "github.com/elithrar/station"
)

func main() {
    r := http.NewServeMux()

    // IndexHandler is just a func(w http.ResponseWriter, r *http.Request) here.
    r.HandleFunc("/", IndexHandler)

    // Just wrap any handler with station.Cache and it'll set HTTP caching 
    // headers to one week from now. It makes a ton of sense to set long-lived
    // headers on our static content, so let's do that.
    // PS: Like before, call goji.Use(station.Cache()) if you're using Goji.
    r.Handle("/static/", station.Cache(station.Serve("/Users/matt/Desktop/static",
    station.DirList()))

    http.ListenAndServe(":8000", r)
}
```

You can otherwise call `station.Cache(station.MaxAge(3600))` to change the cache
duration (defined in seconds) to whatever you'd like.

If you still have questions about how to use station, raise an issue with 
your code (as it stands) and the error message (if any). PRs for bugs and new 
features are welcomed, with API stability the #1 priority.

# License

BSD licensed. See the LICENSE file for details.

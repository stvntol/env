package env

import (
	"context"
	"net/http"
	"path"
	"strconv"
	"strings"
)

/* Shift Path Routing */

// Router is a Handler which servers an http.Handler based on a http.Request.
type Router interface {
	Handler
	Handler(r *http.Request) (h http.Handler)
}

// router implements Router.
type router struct {
	E *Env
	R RouterFunc
}

// RouterFunc is a function which picks an http.Handler based off the value of
// a 'head' string (section or the r.URL.Path). It takes an *Env so that the
// returned http.Handlers can be Handlers.
type RouterFunc func(e *Env, head string) (h http.Handler)

// Env makes router implement Handler
func (rt router) Env() *Env {
	return rt.E
}

// ServeHTTP makes router implement http.Handler
func (rt router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.Handler(r).ServeHTTP(w, r)
}

// Handler allows for the retrieval of an http.Handler based on the
// http.Request. Each successive Router steps through http.Request.URL.Path as
// it routes.
// For example for a request to /users/{user-id}/orders the first router only
// considers the value 'users' to decided which handler to return. If it
// returns another Router that one will only consider the value {user-id}. If
// that returns another Router that router will only consider 'orders' to
// decided which handler to return.
// Based on the routing method described here:
// https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html
func (rt router) Handler(r *http.Request) (h http.Handler) {
	head, _ := shiftPathDepth(r)
	return rt.R(rt.E, head)
}

func shiftPathDepth(r *http.Request) (head string, depth int) {

	path := r.URL.Path
	ctx := r.Context()

	depth = PathDepthFromContext(ctx)

	head, _ = ShiftPath(path, depth)
	depth++

	*r = *r.WithContext(WithPathDepth(ctx, depth))

	return head, depth

}

// ShiftPath cleans and selects a part of a path. If the path was /a/b/c/d and
// ShiftPath was passed depth = 2. It would return head = c, tail = /d. If
// depth = 3 it would return head = d, tail = / and any depth greater than 3
// will return head = '', tail = /
func ShiftPath(p string, depth int) (head, tail string) {
	p = path.Clean("/" + p)

	for d := 0; d <= depth; d++ {
		i := strings.Index(p[1:], "/") + 1
		if i <= 0 {
			if depth > d {
				return "", "/"
			}
			return p[1:], "/"
		}
		head = p[1:i]
		p = p[i:]
	}

	return head, p
}

/* Context */

type contextPathDepthType struct{}

var contextPathDepthKey = &contextPathDepthType{}

const requestPathDepthHeader = "request-path-depth"

// WithPathDepth puts the request URL path depth into the current context.
func WithPathDepth(ctx context.Context, depth int) context.Context {
	return context.WithValue(ctx, contextPathDepthKey, depth)
}

// PathDepthFromContext returns the path depth from context. If there is no
// depth in the current context it returns 0.
func PathDepthFromContext(ctx context.Context) int {
	v := ctx.Value(contextPathDepthKey)
	if v == nil {
		return 0
	}
	return v.(int)
}

// PathDepthHandler is a middleware which places the routing path depth into
// the request context if it's found in the request header.
func PathDepthHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		depth := r.Header.Get(requestPathDepthHeader)

		d, err := strconv.Atoi(depth)
		if err == nil {
			r = r.WithContext(WithPathDepth(r.Context(), d))
			r.Header.Del(requestPathDepthHeader)
		}

		h.ServeHTTP(w, r)
	})

}

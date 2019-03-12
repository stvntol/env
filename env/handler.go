package env

import (
	"context"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
)

// Error is an error which provides access to a HTTP status code and an
// optional message for the HTTP response.
type Error interface {
	error
	Status() int
	Message() string
}

// StatusError is an error which implements the Error interface.
type StatusError struct {
	Code int
	Err  error
	Msg  string
}

// Error allows StatusError to implement the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
}

// Status allows StatusError to implement the Error interface.
func (se StatusError) Status() int {
	return se.Code
}

// Message allows StatusError to implement Error. It returns an
// alternative message intended for the HTTP response.
func (se StatusError) Message() string {
	if se.Msg == "" {
		return se.Error()
	}
	return se.Msg
}

// HandlerFunc is a function to be used as a special HTTP handler that takes an
// Env and returns an error.
type HandlerFunc func(e *Env, w http.ResponseWriter, r *http.Request) error

// Handler interface is an http.Handler that also provides access to an *Env
type Handler interface {
	Env() *Env
	http.Handler
}

// handler takes a configured Env and a HandlerFunc.
type handler struct {
	E *Env
	H HandlerFunc
}

// Env makes handler implement Handler
func (h handler) Env() *Env {
	return h.E
}

// ServeHTTP allows the handler type to satisfy http.Handler.
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.Env(), w, r)
	if err != nil {
		h.Env().errorHandle(err, h.Env(), w, r)
	}
}

// ErrorHandler handles errors in handler's ServeHTTP method
type ErrorHandler func(error, *Env, http.ResponseWriter, *http.Request)

// DefaultErrorHandler handles errors in handler's ServeHTTP method unless
// another handler was specified when the Env was created.
func DefaultErrorHandler(err error, env *Env, w http.ResponseWriter, r *http.Request) {
	switch e := err.(type) {
	case StatusError:
		// We can retrieve the status here and write out a specific
		// HTTP status code.
		log.Printf("HTTP %d - %s %s", e.Status(), e, r.URL.Path)
		http.Error(w, e.Message(), e.Status())
	default:
		// Any error types we don't specifically look out for default
		// to serving a HTTP 500
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}

}

/* ServeMux */

// ServeMux is an HTTP request multiplexer that registers Handlers.
type ServeMux struct {
	e   *Env
	mux *http.ServeMux
}

// NewServeMux allocates and returns a new EnvServeMux.
func NewServeMux(env *Env) *ServeMux {
	return &ServeMux{env, http.NewServeMux()}
}

// HandleFunc registers the handler function for a given pattern.
func (esm *ServeMux) HandleFunc(pattern string, hf HandlerFunc) {
	esm.mux.Handle(pattern, handler{esm.e, hf})
}

// Handle registers the handler for a given pattern.
func (esm *ServeMux) Handle(pattern string, handler Handler) {
	esm.mux.Handle(pattern, handler)
}

// ServeHTTP dispatches the request to the handler whose pattern most closely
// matches the request URL.
func (esm *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	esm.mux.ServeHTTP(w, r)
}

// Env makes ServeMux implement Handler
func (esm *ServeMux) Env() *Env {
	return esm.e
}

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

// PathDepthHandler is a middleware which places the routing path depth in to
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

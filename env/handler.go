package env

import (
	"log"
	"net/http"
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
	path := r.URL.Path
	err := h.H(h.Env(), w, r)
	if err != nil {
		switch e := err.(type) {
		case StatusError:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s %s", e.Status(), e, path)
			http.Error(w, e.Message(), e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

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

func (esm *ServeMux) Env() *Env {
	return esm.e
}

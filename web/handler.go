package web

import (
	"github.com/stvntol/dt/env"

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
type HandlerFunc func(e *env.Env, w http.ResponseWriter, r *http.Request) error

// Handler takes a configured Env and a HandlerFunc.
type Handler struct {
	*env.Env
	H HandlerFunc
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.Env, w, r)
	if err != nil {
		switch e := err.(type) {
		case StatusError:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Message(), e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

// EnvServeMux is an HTTP request multiplexer that registers Handlers.
type EnvServeMux struct {
	*env.Env
	mux *http.ServeMux
}

// NewEnvServeMux allocates and returns a new EnvServeMux.
func NewEnvServeMux(env *env.Env) *EnvServeMux {
	return &EnvServeMux{env, http.NewServeMux()}
}

// HandleFunc registers the handler function for a given pattern.
func (esm *EnvServeMux) HandleFunc(pattern string, handler HandlerFunc) {
	esm.mux.Handle(pattern, Handler{esm.Env, handler})
}

// Handle registers the handler for a given pattern.
func (esm *EnvServeMux) Handle(pattern string, handler Handler) {
	esm.mux.Handle(pattern, handler)
}

// ServeHTTP dispatches the request to the handler whose pattern most closely
// matches the request URL.
func (esm *EnvServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	esm.mux.ServeHTTP(w, r)
}

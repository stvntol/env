package env

import (
	"log"
	"net/http"
)

// Error is an error which provides access to a HTTP status code and an
// alternative message for the HTTP response.
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

// ServeHTTP makes handler implement http.Handler
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

/* Swap Env Handler */

// SwapCondition is a function used as a test for SwapEnvHandler and
// SwapEnvRouter. It may return a new Env or be passed down to further
// handlers.
type SwapCondition func(e *Env, r *http.Request) (*Env, error)

// SwapEnvHandler returns a Handler made from the HandlerFunc passed to it if
// the SwapCondition passed to it returns a nil error. If the SwapCondition
// returns an error that error will be handled by the ErrorHandler of the Env
// passed to it and not the one return by SwapCondtion.
func SwapEnvHandler(e *Env, con SwapCondition, h HandlerFunc) Handler {
	return e.HandlerFunc(func(e *Env, w http.ResponseWriter, r *http.Request) error {

		newEnv, err := con(e, r)
		if err != nil {
			return err
		}

		newEnv.HandlerFunc(h).ServeHTTP(w, r)
		return nil
	})

}

// SwapEnvRouter returns a Handler made from the RouterFunc passed to it if
// the SwapCondition passed to it returns a nil error. If the SwapCondition
// returns an error that error will be handled by the ErrorHandler of the Env
// passed to it and not the one return by SwapCondtion.
func SwapEnvRouter(e *Env, con SwapCondition, routerFn RouterFunc) Handler {
	return e.HandlerFunc(func(e *Env, w http.ResponseWriter, r *http.Request) error {

		newEnv, err := con(e, r)
		if err != nil {
			return err
		}

		newEnv.RouterFunc(routerFn).ServeHTTP(w, r)
		return nil
	})

}

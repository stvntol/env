package env

// Env holds the application wide environment data.
type Env struct {
	value       interface{}
	errorHandle ErrorHandler
}

// HandlerFunc turns a HandlerFunc into a Handler
func (e *Env) HandlerFunc(fn HandlerFunc) handler {
	return handler{e, fn}
}

// RouterFunc turns a RouterFunc into a Router
func (e *Env) RouterFunc(fn RouterFunc) router {
	return router{e, fn}
}

// Value returns whatever environment data was stored when the Env was created.
func (e *Env) Value() interface{} {
	return e.value
}

// ErrorHandler returns the environement's error handler
func (e *Env) ErrorHandler() ErrorHandler {
	return e.errorHandle
}

// NewEnv returns a pointer to a new Env.  If no ErrorHandler is provided
// DefaultErrorHandler is used.
func NewEnv(val interface{}, eh ErrorHandler) *Env {
	if eh == nil {
		return &Env{val, DefaultErrorHandler}
	}
	return &Env{val, eh}
}

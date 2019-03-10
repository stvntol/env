package env

// Env holds the application wide environment configuration.
type Env struct {
	Host string
}

func (e *Env) HandlerFunc(fn HandlerFunc) Handler {
	return handler{e, fn}
}

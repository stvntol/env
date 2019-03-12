package main

import (
	"errors"
	"github.com/stvntol/dt/env"
	"net/http"
)

type Requester struct {
	id       string
	Username string
	token    string
	d        *DataSource
}

func envRequester(e *env.Env) *Requester {
	return e.Value().(*Requester)
}

func RequesterAuth(e *env.Env, r *http.Request) (*env.Env, error) {

	//Do auth stuff
	r.ParseForm()
	username := r.FormValue("username")
	id := r.RemoteAddr

	// failure return error
	if username != "admin" {
		return nil, env.StatusError{
			Code: http.StatusUnauthorized,
			Err:  errors.New(http.StatusText(http.StatusUnauthorized)),
		}
	}

	cr := &Requester{
		id:       id,
		Username: username,
		token:    "token1234",
		d:        envData(e),
	}

	return SwapEnvVal(e, cr), nil
}

func SwapEnvVal(e *env.Env, val interface{}) *env.Env {
	return env.NewEnv(val, e.ErrorHandler())
}

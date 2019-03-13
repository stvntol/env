package main

import (
	"errors"
	"github.com/stvntol/dt/env"
	"log"
	"net/http"
	"strings"
)

func main() {

	e := env.NewEnv(&DataSource{}, nil)

	index := e.RouterFunc(IndexRouter)

	log.Printf("Server start")
	err := http.ListenAndServe(":7357", index)
	log.Printf("ERROR:Server Stopped: %s", err)
}

func envData(e *env.Env) *DataSource {
	return e.Value().(*DataSource)
}

type DataSource struct {
}

func (d *DataSource) String() string {
	return "This is a data source"
}

type dumb string

func (d dumb) String() string {
	return strings.Title(string(d))
}
func (d dumb) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(d.String()))
}
func (d dumb) F(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte(d.String()))
	return nil
}

func StatusNotFoundHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	return env.StatusError{
		Code: http.StatusNotFound,
		Err:  errors.New(http.StatusText(http.StatusNotFound)),
	}
}

package main

import (
	"github.com/stvntol/dt/env"
	"github.com/stvntol/dt/web"
	"log"
	"net/http"
)

func main() {

	env := &env.Env{"This isn't really a host name"}

	mux := web.NewEnvServeMux(env)

	mux.HandleFunc("/", Index)

	log.Printf("Server start")
	err := http.ListenAndServe(":7357", mux)
	log.Printf("ERROR:Server Stopped: %s", err)
}

func Index(e *env.Env, w http.ResponseWriter, r *http.Request) error {

	w.Write([]byte(e.Host))
	return nil

}

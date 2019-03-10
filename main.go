package main

import (
	"errors"
	"github.com/stvntol/dt/env"
	"log"
	"net/http"
	"path"
	"strings"
)

func main() {

	e := &env.Env{"This isn't really a host name"}

	index := e.HandlerFunc(Index)

	log.Printf("Server start")
	err := http.ListenAndServe(":7357", index)
	log.Printf("ERROR:Server Stopped: %s", err)
}

func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func Index(e *env.Env, w http.ResponseWriter, r *http.Request) error {

	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)

	switch head {
	case "favicon.ico":
		http.NotFound(w, r)

	case "accounts":
		w.Write([]byte("Accounts"))

	case "users":
		// TODO figure out way to preserve full path in errors after ShiftPath.
		e.HandlerFunc(UsersHandler).ServeHTTP(w, r)

	case "restaurants":
		w.Write([]byte("Restaurants"))

	case "tables":
		w.Write([]byte("Tables"))

	default:
		w.Write([]byte(e.Host))
	}

	return nil

}

func UsersHandler(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	var username string
	username, r.URL.Path = ShiftPath(r.URL.Path)

	if username == "" {
		// handle
		//
		// 	GET
		// 		get list of users/usernames
		// 	POST
		// 		register a new user. returns a session token
		//
		w.Write([]byte("Users"))
		return nil
	}

	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)

	switch head {

	case "avatar":
	case "status":
	case "bookmarks":
	case "friends":
	case "password":
	case "ratings":

	case "session":
	case "email-verification":

	case "error":
		return env.StatusError{
			500,
			errors.New("this was an intened error shh!"),
			"Something went wrong " + username + " =(",
		}

	default:
		w.Write([]byte("User: " + username + " " + head))
	}

	return nil
}

type UserHandler struct {
	ID       string
	Username string
	e        *env.Env
}

func (u *UserHandler) Message(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	msg := "Generic Message: " + e.Host + " " + u.Username
	w.Write([]byte(msg))
	return nil
}

func (u *UserHandler) Greeting(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	msg := "Hello " + u.Username
	w.Write([]byte(msg))
	return nil
}

func verifyMiddleware(next env.Handler) env.Handler {

	return next.Env().HandlerFunc(
		func(e *env.Env, w http.ResponseWriter, r *http.Request) error {

			// Do verification stuff
			r.ParseForm()
			username := r.FormValue("username")
			// id := r.RemoteAddr

			// failure return error
			if username == "" {
				return env.StatusError{
					Code: http.StatusUnauthorized,
					Err:  errors.New(http.StatusText(http.StatusUnauthorized)),
				}
			}

			// user := UserHandler{
			// 	id,
			// 	username,
			// 	next.Env(),
			// }

			next.ServeHTTP(w, r)
			return nil
		})
}

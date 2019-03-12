package main

import (
	"errors"
	"github.com/stvntol/dt/env"
	"log"
	"net/http"
)

func main() {

	e := env.NewEnv("This isn't really a host name", nil)

	index := e.RouterFunc(IndexRouter)

	log.Printf("Server start")
	err := http.ListenAndServe(":7357", index)
	log.Printf("ERROR:Server Stopped: %s", err)
}

func envData(e *env.Env) string {
	return e.Value().(string)
}

func IndexRouter(e *env.Env, head string) http.Handler {

	switch head {
	case "favicon.ico":
		return http.HandlerFunc(http.NotFound)

	case "accounts":
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Accounts"))
		})

	case "users":
		return e.RouterFunc(UsersRouter)

	case "restaurants":
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Restaurants"))
		})

	case "tables":
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Tables"))
		})

	default:
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(envData(e)))
		})
	}

}

func UsersRouter(e *env.Env, head string) http.Handler {

	if head == "" {
		// handle
		//
		// 	GET
		// 		get list of users/usernames
		// 	POST
		// 		register a new user. returns a session token
		//
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Users"))
		})
	}

	return e.RouterFunc(UsernameRouter(head))

}

func UsernameRouter(username string) env.RouterFunc {
	return func(e *env.Env, head string) http.Handler {
		switch head {
		//		case "avatar":
		//		case "status":
		//		case "bookmarks":
		//		case "friends":
		//		case "password":
		//		case "ratings":
		//
		//		case "session":
		//		case "email-verification":

		case "error":
			return e.HandlerFunc(func(e *env.Env, w http.ResponseWriter, r *http.Request) error {
				return env.StatusError{
					Code: 500,
					Err:  errors.New("this was an intened error shh! error at head: " + head),
					Msg:  "Something went wrong " + username + " =(",
				}
			})

		default:
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("User: " + username + " " + head))
			})
		}

	}
}

type UserHandler struct {
	ID       string
	Username string
	e        *env.Env
}

func (u *UserHandler) Message(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	msg := "Generic Message: " + envData(e) + " " + u.Username
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

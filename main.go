package main

import (
	"errors"
	"github.com/stvntol/dt/env"
	"log"
	"net/http"
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

func IndexRouter(e *env.Env, head string) http.Handler {

	switch head {
	case "favicon.ico":
		return http.HandlerFunc(http.NotFound)

	case "accounts":
		return dumb("Accounts")

	case "health":
		return dumb("Status Ok")

	case "restaurants":
		return dumb("Restaurants")

	case "tables":
		return env.SwapHandlerFunc(e, RequesterAuth, dumb("Tables").F)

	case "users":
		return env.SwapRouterFunc(e, RequesterAuth, UsersRouter)

	default:
		return dumb(envData(e).String())
	}

}

type dumb string

func (d dumb) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(d))
}
func (d dumb) F(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte(d))
	return nil
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
		return dumb("Users")
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
	msg := "Generic Message: " + envData(e).String() + " " + u.Username
	w.Write([]byte(msg))
	return nil
}

func (u *UserHandler) Greeting(e *env.Env, w http.ResponseWriter, r *http.Request) error {
	msg := "Hello " + u.Username
	w.Write([]byte(msg))
	return nil
}

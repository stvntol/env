package main

import (
	"fmt"
	"github.com/stvntol/dt/env"
	"net/http"
)

func IndexRouter(e *env.Env, head string) http.Handler {

	switch head {

	case "accounts":
		return e.RouterFunc(AccountsRouter)

	case "favicon.ico":
		return http.HandlerFunc(http.NotFound)

	case "health":
		return dumb("Status Ok")

	case "view":
		return e.RouterFunc(ViewRouter)

	/**** Require Authorization ****/

	case "restaurants":
		return env.SwapRouterFunc(e, RequesterAuth, RestaurantsRouter)

	case "tables":
		return env.SwapRouterFunc(e, RequesterAuth, TablesRouter)

	case "users":
		return env.SwapRouterFunc(e, RequesterAuth, UsersRouter)

	case "ws":
		return env.SwapHandlerFunc(e, RequesterAuth, dumb("WS").F)

	default:
		// This should match the "view" case to match on short table view URLs
		return env.SwapHandlerFunc(e, RequesterAuth, StatusNotFoundHandler)
	}

}

func AccountsRouter(e *env.Env, head string) http.Handler {

	switch head {

	case "":
		return dumb(fmt.Sprintf("New Account"))

	case "email-verification":
		return dumb(head)

	case "password":
		return dumb(fmt.Sprintf("New Account"))

	default:
		return e.RouterFunc(AccountsUsernameRouter(head))
	}
}

func AccountsUsernameRouter(username string) env.RouterFunc {
	return func(e *env.Env, head string) http.Handler {
		switch head {

		case "":
			return dumb(fmt.Sprintf("New Account"))

		case "email-verification":
			return dumb(head)

		case "password":
			return dumb(fmt.Sprintf("New Account"))

		default:
			return e.HandlerFunc(StatusNotFoundHandler)
		}
	}
}

func RestaurantsRouter(e *env.Env, head string) http.Handler {

	if head == "" {
		return dumb("Restaurants")
	}

	// return e.HandlerFunc(RestaurantCodeHandler(head))
	return dumb(fmt.Sprintf("Restaurant code %s", head))
}

func TablesRouter(e *env.Env, head string) http.Handler {

	if head == "" {
		return dumb("Tables")
	}

	return e.RouterFunc(TableIDRouter(head))
}

func TableIDRouter(tableID string) env.RouterFunc {
	return func(e *env.Env, head string) http.Handler {
		switch head {
		case "":
			return dumb(fmt.Sprintf("Table ID %s", tableID, head))
		case "chat":
			return dumb(fmt.Sprintf("Table ID %s's %s", tableID, head))
		case "invite":
			return dumb(fmt.Sprintf("Table ID %s's %s", tableID, head))

		default:
			return e.HandlerFunc(StatusNotFoundHandler)
		}
	}
}

func UsersRouter(e *env.Env, head string) http.Handler {

	if head == "" {
		return dumb("Users")
	}

	return e.RouterFunc(UsernameRouter(head))

}

func UsernameRouter(username string) env.RouterFunc {
	return func(e *env.Env, head string) http.Handler {
		switch head {

		case "avatar":
			return dumb(fmt.Sprintf("%s's %s", username, head))

		case "bookmarks":
			return dumb(fmt.Sprintf("%s's %s", username, head))

		case "friends":
			return dumb(fmt.Sprintf("%s's %s", username, head))

		case "password":
			return dumb(fmt.Sprintf("%s's %s", username, head))

		case "ratings":
			return dumb(fmt.Sprintf("%s's %s", username, head))

		case "session":
			return dumb(fmt.Sprintf("%s's %s", username, head))

		case "status":
			return dumb(fmt.Sprintf("%s's %s", username, head))

		default:
			return e.HandlerFunc(StatusNotFoundHandler)

		}

	}
}

func ViewRouter(e *env.Env, tableID string) http.Handler {

	if tableID == "" {
		return e.HandlerFunc(StatusNotFoundHandler)
	}

	// return e.HanderFunc(ViewTableHandler(tableID))
	return dumb(fmt.Sprintf("View table %s", tableID))
}

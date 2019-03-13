package main

import (
	_ "fmt"
	"net/http"
)

func UsersHandler(req *Requester, w http.ResponseWriter, r *http.Request) error {
	dumb("Users").ServeHTTP(w, r)
	return nil
}

type Username string

func (u Username) AvatarHandler(req *Requester, w http.ResponseWriter, r *http.Request) error {
	dumb("avatar for "+u).ServeHTTP(w, r)
	return nil
}

func (u Username) BookmarksHandler(req *Requester, w http.ResponseWriter, r *http.Request) error {
	dumb("bookmarks for "+u).ServeHTTP(w, r)
	return nil
}
func (u Username) FriendsHandler(req *Requester, w http.ResponseWriter, r *http.Request) error {
	dumb("friends for "+u).ServeHTTP(w, r)
	return nil
}
func (u Username) PasswordHandler(req *Requester, w http.ResponseWriter, r *http.Request) error {
	dumb("password for "+u).ServeHTTP(w, r)
	return nil
}
func (u Username) RatingsHandler(req *Requester, w http.ResponseWriter, r *http.Request) error {
	dumb("ratings for "+u).ServeHTTP(w, r)
	return nil
}
func (u Username) SessionHandler(req *Requester, w http.ResponseWriter, r *http.Request) error {
	dumb("session for "+u).ServeHTTP(w, r)
	return nil
}
func (u Username) StatusHandler(req *Requester, w http.ResponseWriter, r *http.Request) error {
	dumb("status for "+u).ServeHTTP(w, r)
	return nil
}

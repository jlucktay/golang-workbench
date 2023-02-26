// Package server defines internal behaviour.
package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	root "go.jlucktay.dev/golang-workbench/go_rest_api/pkg"
)

type userRouter struct {
	userService root.UserService
}

// NewUserRouter sets up a Router with some functions to handle various routes.
func NewUserRouter(u root.UserService, router *mux.Router) *mux.Router {
	userRouter := userRouter{u}

	router.HandleFunc("/", userRouter.createUserHandler).Methods("PUT")
	router.HandleFunc("/profile", validate(userRouter.profileHandler)).Methods("GET")
	router.HandleFunc("/{username}", userRouter.getUserHandler).Methods("GET")
	router.HandleFunc("/login", userRouter.loginHandler).Methods("POST")

	return router
}

func (s *userRouter) createUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err := decodeUser(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = s.userService.CreateUser(&user)

	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSON(w, http.StatusOK, err)
}

func (s *userRouter) profileHandler(w http.ResponseWriter, r *http.Request) {
	claim, ok := r.Context().Value(contextKeyAuthtoken).(claims)

	if !ok {
		Error(w, http.StatusBadRequest, "no context")
		return
	}

	username := claim.Username
	user, err := s.userService.GetUserByUsername(username)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	JSON(w, http.StatusOK, user)
}

func (s *userRouter) getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println(vars)
	username := vars["username"]

	user, err := s.userService.GetUserByUsername(username)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	JSON(w, http.StatusOK, user)
}

func (s *userRouter) loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("loginHandler")
	credentials, err := decodeCredentials(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var user root.User
	user, err = s.userService.Login(credentials)

	if err == nil {
		cookie := newAuthCookie(user)
		JSONWithCookie(w, http.StatusOK, user, cookie)
	} else {
		Error(w, http.StatusInternalServerError, "Incorrect password")
	}
}

func decodeUser(r *http.Request) (root.User, error) {
	var u root.User

	if r.Body == nil {
		return u, errors.New("no request body")
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)

	return u, err
}

func decodeCredentials(r *http.Request) (root.Credentials, error) {
	var c root.Credentials

	if r.Body == nil {
		return c, errors.New("no request body")
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&c)

	return c, err
}

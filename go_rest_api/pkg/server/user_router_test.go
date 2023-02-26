package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"

	root "go.jlucktay.dev/golang-workbench/go_rest_api/pkg"
	"go.jlucktay.dev/golang-workbench/go_rest_api/pkg/mock"
)

// createUserHandler tests
func Test_UserRouter_createUserHandler(t *testing.T) {
	t.Run("happy path", createUserHandlerShouldPassUserObjectToUserServiceCreateUser)
	t.Run("invalid payload", createUserHandlerShouldReturnStatusBadRequestIfPayloadIsInvalid)
	t.Run("internal error", createUserHandlerShouldReturnStatusInternalServerErrorIfUserServiceReturnsError)
}

func createUserHandlerShouldPassUserObjectToUserServiceCreateUser(t *testing.T) {
	// Arrange
	us := mock.UserService{}
	testMux := NewUserRouter(&us, mux.NewRouter())
	var result *root.User
	us.CreateUserFn = func(u *root.User) error {
		result = u
		return nil
	}

	testUsername := "test_username"
	testPassword := "test_password"

	values := map[string]string{"username": testUsername, "password": testPassword}
	jsonValue, _ := json.Marshal(values)
	payload := bytes.NewBuffer(jsonValue)

	// Act
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/", payload)
	r.Header.Set("Content-Type", "application/json")
	testMux.ServeHTTP(w, r)

	// Assert
	if !us.CreateUserInvoked {
		t.Fatal("expected CreateUser() to be invoked")
	}
	if result.Username != testUsername {
		t.Fatalf("expected username to be: `%s`, got: `%s`", testUsername, result.Username)
	}
	if result.Password != testPassword {
		t.Fatalf("expected username to be: `%s`, got: `%s`", testPassword, result.Password)
	}
}

func createUserHandlerShouldReturnStatusBadRequestIfPayloadIsInvalid(t *testing.T) {
	// Arrange
	us := mock.UserService{}
	testMux := NewUserRouter(&us, mux.NewRouter())
	us.CreateUserFn = func(u *root.User) error {
		return nil
	}

	// Act
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/", nil)
	r.Header.Set("Content-Type", "application/json")
	testMux.ServeHTTP(w, r)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Fatal("expected: http.StatusBadRequest, got: %i", w.Code)
	}
}

func createUserHandlerShouldReturnStatusInternalServerErrorIfUserServiceReturnsError(t *testing.T) {
	// Arrange
	us := mock.UserService{}
	testMux := NewUserRouter(&us, mux.NewRouter())
	us.CreateUserFn = func(u *root.User) error {
		return errors.New("user service error")
	}

	values := map[string]string{"username": "", "password": ""}
	jsonValue, _ := json.Marshal(values)
	payload := bytes.NewBuffer(jsonValue)

	// Act
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/", payload)
	r.Header.Set("Content-Type", "application/json")
	testMux.ServeHTTP(w, r)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Fatal("expected: http.StatusInternalServerError, got: %i", w.Code)
	}
}

// profileHandler tests
func Test_UserRouter_profileHandler(t *testing.T) {
	t.Run("happy path", profileHandlerShouldReturnUserFromContext)
	t.Run("no context", profileHandlerShouldReturnStatusBadRequestIfNoAuthContext)
	t.Run("user not found", profileHandlerShouldReturnStatusNotFoundIfNoUserFound)
}

func profileHandlerShouldReturnUserFromContext(t *testing.T) {
	// Arrange
	us := mock.UserService{}
	testMux := NewUserRouter(&us, mux.NewRouter())
	var result string
	us.GetUserByUsernameFn = func(username string) (root.User, error) {
		result = username
		return root.User{}, nil
	}

	testUsername := "test_username"
	testUser := root.User{Username: testUsername}

	// Act
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/profile", nil)
	testCookie := newAuthCookie(testUser)
	r.AddCookie(&testCookie)
	ctx := context.WithValue(r.Context(), contextKeyAuthtoken, claims{testUsername, jwt.StandardClaims{}})
	testMux.ServeHTTP(w, r.WithContext(ctx))

	// Assert
	if !us.GetUserByUsernameInvoked {
		t.Fatal("expected GetUserByUsername() to be invoked")
	}
	if result != testUsername {
		t.Fatalf("expected username to be: `%s`, got: `%s`", testUsername, result)
	}
}

func profileHandlerShouldReturnStatusBadRequestIfNoAuthContext(t *testing.T) {
	// Arrange
	us := mock.UserService{}
	testMux := NewUserRouter(&us, mux.NewRouter())

	// Act
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/profile", nil)
	testMux.ServeHTTP(w, r)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected StatusUnauthorized, got: %d", w.Code)
	}
}

func profileHandlerShouldReturnStatusNotFoundIfNoUserFound(t *testing.T) {
	// Arrange
	us := mock.UserService{}
	testMux := NewUserRouter(&us, mux.NewRouter())
	us.GetUserByUsernameFn = func(username string) (root.User, error) {
		return root.User{}, errors.New("user service error")
	}
	testUsername := "test_username"
	testUser := root.User{Username: testUsername}

	// Act
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/profile", nil)
	testCookie := newAuthCookie(testUser)
	r.AddCookie(&testCookie)
	ctx := context.WithValue(r.Context(), contextKeyAuthtoken, claims{testUsername, jwt.StandardClaims{}})
	testMux.ServeHTTP(w, r.WithContext(ctx))

	// Assert
	if !us.GetUserByUsernameInvoked {
		t.Fatal("expected GetUserByUsername() to be invoked")
	}
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected: StatusNotFound, got: %d", w.Code)
	}
}

// getUserHandler tests
func Test_UserRouter_getUserHandler(t *testing.T) {
	t.Run("happy path", getUserHandlerShouldCallGetUserByUsernameWithUsernameFromQuerystring)
	t.Run("no user found", getUserHandlerShouldReturnStatusNotFoundIfNoUserFound)
}

func getUserHandlerShouldCallGetUserByUsernameWithUsernameFromQuerystring(t *testing.T) {
	// Arrange
	us := mock.UserService{}
	testMux := NewUserRouter(&us, mux.NewRouter())
	var result string
	us.GetUserByUsernameFn = func(username string) (root.User, error) {
		result = username
		return root.User{}, nil
	}

	testUsername := "test_username"

	// Act
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/"+testUsername, nil)
	testMux.ServeHTTP(w, r)

	// Assert
	if !us.GetUserByUsernameInvoked {
		t.Fatal("expected GetUserByUsername() to be invoked")
	}
	if result != testUsername {
		t.Fatalf("expected username to be: `%s`, got: `%s`", testUsername, result)
	}
}

func getUserHandlerShouldReturnStatusNotFoundIfNoUserFound(t *testing.T) {
	// Arrange
	us := mock.UserService{}
	testMux := NewUserRouter(&us, mux.NewRouter())
	us.GetUserByUsernameFn = func(username string) (root.User, error) {
		return root.User{}, errors.New("user service error")
	}

	testUsername := "test_username"

	// Act
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/"+testUsername, nil)
	testMux.ServeHTTP(w, r)

	// Assert
	if !us.GetUserByUsernameInvoked {
		t.Fatal("expected GetUserByUsername() to be invoked")
	}
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected: StatusNotFound, got: %d", w.Code)
	}
}

// gHandler tests
func Test_UserRouter_loginHandler(t *testing.T) {
	fmt.Println("loginHandler tests")
	t.Run("happy path", loginHandlerShouldProvideNewAuthCookieIfUserServiceReturnsAUser)
	// t.Run("no user found", getUserHandlerShouldReturnStatusNotFoundIfNoUserFound)
}

func loginHandlerShouldProvideNewAuthCookieIfUserServiceReturnsAUser(t *testing.T) {
	// Arrange
	us := mock.UserService{}
	testMux := NewUserRouter(&us, mux.NewRouter())
	us.LoginFn = func(credentials root.Credentials) (root.User, error) {
		return root.User{}, nil
	}

	testUsername := "test_username"
	testPassword := "test_password"

	values := map[string]string{"username": testUsername, "password": testPassword}
	jsonValue, _ := json.Marshal(values)
	payload := bytes.NewBuffer(jsonValue)

	// Act
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/login", payload)
	testMux.ServeHTTP(w, r)

	// Assert
	if !us.LoginInvoked {
		t.Fatal("expected Login() to be invoked")
	}

	request := &http.Request{Header: http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}}
	cookie, err := request.Cookie("Auth")
	if err != nil || cookie == nil {
		panic("Expected Cookie named 'Auth'")
	}
}

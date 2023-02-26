// Package server defines internal behaviour.
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	root "go.jlucktay.dev/golang-workbench/go_rest_api/pkg"
)

type claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type contextKey string

func (c contextKey) String() string {
	return "mypackage context key " + string(c)
}

var contextKeyAuthtoken = contextKey("auth-token")

func newAuthCookie(user root.User) http.Cookie {
	expireTime := time.Now().Add(time.Hour * 1)
	c := claims{
		user.Username,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "localhost!",
		},
	}

	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte("secret"))

	cookie := http.Cookie{
		Name:     "Auth",
		Value:    token,
		Expires:  expireTime,
		HttpOnly: true,
	}
	return cookie
}

func validate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("Auth")
		if err != nil {
			Error(res, http.StatusUnauthorized, "No authorization cookie")
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected siging method")
			}

			return []byte("secret"), nil
		})
		if err != nil {
			Error(res, http.StatusUnauthorized, "Invalid token")
			return
		}

		if claims, ok := token.Claims.(*claims); ok && token.Valid {
			ctx := context.WithValue(req.Context(), contextKeyAuthtoken, *claims)
			next(res, req.WithContext(ctx))
		} else {
			Error(res, http.StatusUnauthorized, "Unauthorized")
			return
		}
	})
}

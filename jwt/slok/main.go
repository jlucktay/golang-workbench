package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	mySigningKey = "WOW,MuchShibe,ToDogge"
)

func main() {
	createdToken, err := exampleNew([]byte(mySigningKey))
	if err != nil {
		fmt.Println("Creating token failed")
	}
	exampleParse(createdToken, mySigningKey)
}

func exampleNew(mySigningKey []byte) (string, error) {
	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(mySigningKey)
	return tokenString, err
}

func exampleParse(myToken, myKey string) {
	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(myKey), nil
	})

	if err == nil && token.Valid {
		fmt.Println("Your token is valid.  I like your style.")
		fmt.Printf("%+v\n", token)
	} else {
		fmt.Println("This token is terrible!  I cannot accept this.")
	}
}

// Package root defines some fundamental types.
package root

// Credentials represent some credentials to log in with.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

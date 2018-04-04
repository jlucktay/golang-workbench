// Package root defines some fundamental types.
package root

// User represents a single user.
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserService defines an interface for a user service.
type UserService interface {
	CreateUser(u *User) error
	GetUserByUsername(username string) (User, error)
	Login(c Credentials) (User, error)
}

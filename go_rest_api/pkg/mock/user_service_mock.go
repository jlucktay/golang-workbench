// Package mock provides some mocks to test with.
package mock

import (
	"github.com/jlucktay/golang-workbench/go_rest_api/pkg"
)

// UserService has some mock functions.
type UserService struct {
	CreateUserFn      func(u *root.User) error
	CreateUserInvoked bool

	GetUserByUsernameFn      func(username string) (root.User, error)
	GetUserByUsernameInvoked bool

	LoginFn      func(c root.Credentials) (root.User, error)
	LoginInvoked bool
}

// CreateUser flubs user creation.
func (us *UserService) CreateUser(u *root.User) error {
	us.CreateUserInvoked = true
	return us.CreateUserFn(u)
}

// GetUserByUsername flubs fetching of users by name.
func (us *UserService) GetUserByUsername(username string) (root.User, error) {
	us.GetUserByUsernameInvoked = true
	return us.GetUserByUsernameFn(username)
}

// Login flubs logging in.
func (us *UserService) Login(c root.Credentials) (root.User, error) {
	us.LoginInvoked = true
	return us.LoginFn(c)
}

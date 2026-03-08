// Package mongo interfaces with MongoDb for us.
package mongo

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	root "go.jlucktay.dev/golang-workbench/go_rest_api/pkg"
)

// UserService holds a MongoDb collection.
type UserService struct {
	collection *mgo.Collection
}

// NewUserService initialises and returns a new UserService.
func NewUserService(session *mgo.Session) (*UserService, error) {
	collection := session.DB("test").C("user")
	if err := collection.EnsureIndex(userModelIndex()); err != nil {
		return nil, err
	}

	return &UserService{collection}, nil
}

// CreateUser will create a user on the designated service.
func (p *UserService) CreateUser(u *root.User) error {
	user := userModel{Username: u.Username}
	err := user.addSaltedPassword(u.Password)
	if err != nil {
		return err
	}

	return p.collection.Insert(&user)
}

// GetUserByUsername looks up a user by name and returns it.
func (p *UserService) GetUserByUsername(username string) (root.User, error) {
	model := userModel{}
	err := p.collection.Find(bson.M{"username": username}).One(&model)

	return root.User{
			ID:       model.ID.Hex(),
			Username: model.Username,
			Password: "-",
		},
		err
}

// Login will log a user into a given service.
func (p *UserService) Login(c root.Credentials) (root.User, error) {
	model := userModel{}
	errOne := p.collection.Find(bson.M{"username": c.Username}).One(&model)
	if errOne != nil {
		log.Fatal(errOne)
	}

	errPass := model.comparePassword(c.Password)
	if errPass != nil {
		return root.User{}, errPass
	}

	return root.User{
			ID:       model.ID.Hex(),
			Username: model.Username,
			Password: "-",
		},
		errPass
}

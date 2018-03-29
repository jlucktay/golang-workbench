package mongo

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type userModel struct {
	Id           bson.ObjectId `bson:"_id,omitempty"`
	Username     string
	PasswordHash string
	Salt         string
}

func userModelIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

func (u *userModel) comparePassword(password string) error {
	incoming := []byte(password + u.Salt)
	existing := []byte(u.PasswordHash)
	err := bcrypt.CompareHashAndPassword(existing, incoming)
	return err
}

func (u *userModel) addSaltedPassword(password string) error {
	salt := uuid.New().String()
	passwordBytes := []byte(password + salt)
	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hash[:])
	u.Salt = salt

	return nil
}

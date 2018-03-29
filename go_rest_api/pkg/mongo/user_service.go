package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/jlucktay/golang-workbench/go_rest_api/pkg"
)

type UserService struct {
	collection *mgo.Collection
}

func NewUserService(session *Session, dbName string, collectionName string) *UserService {
	collection := session.GetCollection(dbName, collectionName)
	collection.EnsureIndex(userModelIndex())
	return &UserService{collection}
}

func (p *UserService) Create(u *root.User) error {
	user := newUserModel(u)
	return p.collection.Insert(&user)
}

func (p *UserService) GetByUsername(username string) (*root.User, error) {
	model := userModel{}
	err := p.collection.Find(bson.M{"username": username}).One(&model)
	return model.toRootUser(), err
}

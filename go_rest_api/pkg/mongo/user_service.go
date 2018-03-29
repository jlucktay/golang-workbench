package mongo

import (
	"github.com/jlucktay/golang-workbench/go_rest_api/pkg"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserService struct {
	collection *mgo.Collection
}

func NewUserService(session *mgo.Session) *UserService {
	collection := session.DB("test").C("user")
	collection.EnsureIndex(userModelIndex())
	return &UserService{collection}
}

func (p *UserService) CreateUser(u *root.User) error {
	user := userModel{Username: u.Username}
	err := user.addSaltedPassword(u.Password)
	if err != nil {
		return err
	}

	return p.collection.Insert(&user)
}

func (p *UserService) GetUserByUsername(username string) (error, root.User) {
	model := userModel{}
	err := p.collection.Find(bson.M{"username": username}).One(&model)
	return err, root.User{
		Id:       model.Id.Hex(),
		Username: model.Username,
		Password: "-"}
}

func (p *UserService) Login(c root.Credentials) (error, root.User) {
	model := userModel{}
	err := p.collection.Find(bson.M{"username": c.Username}).One(&model)

	err = model.comparePassword(c.Password)
	if err != nil {
		return err, root.User{}
	}

	return err, root.User{
		Id:       model.Id.Hex(),
		Username: model.Username,
		Password: "-"}
}

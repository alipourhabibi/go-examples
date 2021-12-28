package models

import (
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/repo"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/services"
)

type UserRepo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
	}
}

func (u *UserRepo) FindByUserName() (repo.User) {
	return repo.FindByUserName(u.Username)
}

func (u *UserRepo) Save() {
	user := repo.User{}
	hashedPass, err := services.GenerateHash(u.Password)
	if err != nil {
		panic(err)
	}
	user.Username = u.Username
	user.Password = hashedPass
	user.Save()
}

func (u *UserRepo) Exist() bool {
	user := repo.User{}
	user.Username = u.Username
	return user.Exist()
}

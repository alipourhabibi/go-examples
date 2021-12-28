package models

import (
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/repo"
)

type UserRepo struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
	}
}

func (u *UserRepo) FindByUserName() (repo.User) {
	return repo.FindByUserName(u.UserName)
}

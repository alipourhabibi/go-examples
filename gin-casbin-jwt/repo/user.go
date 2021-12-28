package repo

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `json:"username"`
	Password string `json:"password"`
}

func FindByUserName(username string) User {
	var user User
	db := GetDB()
	db.First(&user, "username = ?", username)
	return user
}

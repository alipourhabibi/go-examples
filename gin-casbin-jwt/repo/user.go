package repo

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Password string
}

func AutoMigrateAll() {
	db := GetDB()
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Post{})
}

func FindByUserName(username string) User {
	var user User
	db := GetDB()
	db.First(&user, "username = ?", username)
	return user
}

func (u *User) Save() {
	db := GetDB()
	db.Create(u)
}

func (u *User) Exist() bool {
	var user User
	db := GetDB()
	db.First(&user, "username = ? ", u.Username)
	return user.Username != ""
}

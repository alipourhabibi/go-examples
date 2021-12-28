package repo

import (
	"fmt"

	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/settings"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDB() *gorm.DB {
	host := settings.PostgresqlSettings.Items.Host
	user := settings.PostgresqlSettings.Items.User
	pass := settings.PostgresqlSettings.Items.Password
	dbName := settings.PostgresqlSettings.Items.DBName
	port := settings.PostgresqlSettings.Items.Port
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", host, user, pass, dbName, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

package routers

import (
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/middleware"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/repo"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/routers/api"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	r := gin.Default()
	
	db := repo.GetDB()
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(err)
	}

	r.POST("register", api.Register)
	r.POST("login", api.LogIn)
	r.POST("post", middleware.Authtorized("resource", "write", adapter), api.NewPost)

	r.Use(middleware.AuthMiddleware())
	return r
}
	

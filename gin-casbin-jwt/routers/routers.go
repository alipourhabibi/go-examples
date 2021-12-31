package routers

import (
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/middleware"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/routers/api"

	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	r := gin.Default()
	
	r.POST("register", api.Register)
	r.POST("login", api.LogIn)

	r.POST("post", api.NewPost)
	r.PUT("post/:id", api.UpdatePost)
	r.DELETE("post/:id", api.DeletePost)

	r.Use(middleware.AuthMiddleware())
	return r
}
	

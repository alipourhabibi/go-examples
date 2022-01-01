package routers

import (
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/middleware"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/routers/api"

	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	r := gin.Default()
	
	v1 := r.Group("api/v1")
	{
		v1.POST("register", api.Register)
		v1.POST("login", api.LogIn)
		v1.POST("logout", middleware.AuthMiddleware(), api.LogOut)

		v1.GET("post/:id", api.GetPost)
		v1.POST("post", middleware.AuthMiddleware(), api.NewPost)
		v1.PUT("post/:id", middleware.AuthMiddleware(), api.UpdatePost)
		v1.DELETE("post/:id", middleware.AuthMiddleware(), api.DeletePost)
	}

	return r
}
	

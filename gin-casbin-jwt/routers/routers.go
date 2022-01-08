package routers

import (
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/middleware"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/routers/api"
	_ "github.com/alipourhabibi/go-examples/gin-casbin-jwt/docs"

	"github.com/gin-gonic/gin"
	openapi "github.com/go-openapi/runtime/middleware"
)

func InitRouters() *gin.Engine {
	r := gin.Default()
	
	// make swagger doc
	r.StaticFile("swagger.yaml", "./swagger.yaml")
	opts := openapi.RedocOpts{SpecURL: "./swagger.yaml"}
	sh := openapi.Redoc(opts, nil)
	r.GET("/docs", gin.WrapH(sh))

	v1 := r.Group("api/v1")
	{
		v1.POST("refresh", api.Refresh)

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
	

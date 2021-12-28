package routers

import (
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/routers/api"

	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	r := gin.Default()
	r.POST("login", api.LogIn)
	return r
}
	

package middleware

import (
	"net/http"
	"strings"

	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/services"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/settings"
	gormadapter "github.com/casbin/gorm-adapter/v3"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func (c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		content := strings.Split(authorization, " ")
		if len(content) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
			c.Abort()
			return
		}
		token := content[1]
		secret := settings.AppSettings.Items.JwtAccess
		_, err := services.VerifyJWT(token, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func Authtorized(obj string, act string, adapter *gormadapter.Adapter) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		content := strings.Split(authorization, " ")
		token := content[1]
		username, _ := services.VerifyJWT(token, settings.AppSettings.Items.JwtAccess)

		ok, err := services.Enforce(username, obj, act, adapter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal server error"})
			c.Abort()
			return
		}
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

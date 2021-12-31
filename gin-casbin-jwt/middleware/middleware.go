package middleware

import (
	"net/http"
	"strings"

	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/services"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/settings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func (c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "no Authorization header was given"})
			c.Abort()
			return
		}
		content := strings.Split(authorization, " ")
		if len(content) == 2 {
			token := content[1]
			secret := settings.AppSettings.Items.JwtAccess
			_, err := services.VerifyJWT(token, secret)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
				c.Abort()
				return
			}
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
	}
}

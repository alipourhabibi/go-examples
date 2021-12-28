package api

import (
	"net/http"

	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/models"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/services"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/settings"

	"github.com/gin-gonic/gin"
)

type LogInUser struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func LogIn(c *gin.Context) {
	var user LogInUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON provided"})
		return
	}

	userRepo := models.NewUserRepo()
	userRepo.UserName = user.UserName
	userData := userRepo.FindByUserName()

	if userData.UserName == "" {
		c.JSON(http.StatusNotFound, gin.H{"msg": "User not found"})
		return
	}
	valid := services.ValidatePassword(user.Password, userData.Password)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid Credential"})
		return
	}

	accessToken, err := services.GenerateJWT(userData.UserName, settings.AppSettings.Items.JwtAccess)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal server error"})
		return
	}
	refreshToken , err := services.GenerateJWT(userData.UserName, settings.AppSettings.Items.JwtRefresh)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}

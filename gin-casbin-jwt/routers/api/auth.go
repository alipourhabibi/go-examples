package api

import (
	"net/http"
	"strings"

	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/models"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/repo"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/services"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/settings"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type LogInUser struct {
	// Required: true
	// Example: ali
	Username string `json:"username" binding:"required"`
	// Required: true
	// Example: 1234
	Password string `json:"password" binding:"required"`
}

type RegisterUser struct {
	// Required: true
	// Example: ali
	Username string `json:"username" binding:"required"`
	// Required: true
	// Example: 1234
	Password string `json:"password" binding:"required"`
	// Required: true
	// Example: 1234
	Passwrod2 string `json:"password2" binding:"required"`
}

// swagger:route POST /api/v1/login Login loginUserParameter
//
//	used by user to login to the app
//
//	consumes:
//	- application/json
//
//	produces:
//	- application/json
//
//	schemes: http, https
//	
//	Responses:
//	200: loginSuccess
//	400: responseBadRequest
//	401: responseUnauthorized
//	500: responseInternalServerError
//
// Handler function for login
func LogIn(c *gin.Context) {
	var user LogInUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON provided"})
		return
	}

	userRepo := models.NewUserRepo()
	userRepo.Username = user.Username
	userData := userRepo.FindByUserName()

	if userData.Username == "" {
		c.JSON(http.StatusNotFound, gin.H{"msg": "User not found"})
		return
	}
	valid := services.ValidatePassword(user.Password, userData.Password)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid Credential"})
		return
	}

	td, err := services.CreateTokensAndMetaData(userData.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal server error"})
		return	
	}

	c.JSON(http.StatusOK, gin.H{"access_token": td.AccessToken, "refresh_token": td.RefreshToken})
}

// swagger:route POST /api/v1/register Register registerUserParameter
//
//	used by user to register to the app
//
//	consumes:
//	- application/json
//
//	produces:
//	- application/json
//
//	schemes: http, https
//	
//	Responses:
//	201: responseCreated
//	400: responseBadRequest
//	401: responseUnauthorized
//
// Handler function for register
func Register(c *gin.Context) {
	var user RegisterUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON provided"})
		return
	}

	if user.Password != user.Passwrod2 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Two passwords dont match"})
		return
	}

	u := models.NewUserRepo()
	u.Username = user.Username
	u.Password = user.Password

	exist := u.Exist()
	
	if exist {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "user already exists!"})
		return
	}

	u.Save()


	db := repo.GetDB()
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(err)
	}
	services.AddPolicy(u.Username, u.Username, "update", adapter)
	services.AddPolicy(u.Username, u.Username, "delete", adapter)

	c.JSON(http.StatusCreated, gin.H{"msg": "created"})
}

// swagger:route POST /api/v1/logout LogOut logOutUserParameter
//
//	used by user to logout
//
//	consumes:
//	- application/json
//
//	produces:
//	- application/json
//
//	schemes: http, https
//	
//	Responses:
//	200: responseSuccess
//	400: responseBadRequest
//	401: responseUnauthorized
//
// Handler function to logout
func LogOut(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	content := strings.Split(authorization, " ")
	token := content[1]

	dataMap, err := services.VerifyJWT(token, settings.AppSettings.Items.JwtAccess)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "You are not loged in"})
		return
	}
	claims, ok := dataMap.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	redisClient := repo.GetRedisClient()

	accessUUID, ok := claims["access_uuid"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	username, ok := claims["username"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	refreshUUID := accessUUID + "++" + username

	redisClient.Del(refreshUUID)
	redisClient.Del(accessUUID)
	c.JSON(http.StatusOK, gin.H{"msg": "Success"})

}

// swagger:route POST /api/v1/refresh Refresh refreshTokenParameter
//
//	used by user to refresh token
//
//	consumes:
//	- application/json
//
//	produces:
//	- application/json
//
//	schemes: http, https
//	
//	Responses:
//	200: loginSuccess
//	400: responseBadRequest
//	401: responseUnauthorized
//	500: responseInternalServerError
//
// Handler function to refresh Token
func Refresh(c *gin.Context) {
	token := struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON provided"})
		return
	}

	jwtToken, err := services.VerifyJWT(token.RefreshToken, settings.AppSettings.Items.JwtRefresh)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	// get uuid
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if ok {
		username, usernameOK := claims["username"].(string)		
		if !usernameOK {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
			return
		}
		td, err := services.CreateTokensAndMetaData(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"refresh_token": td.RefreshToken, "access_token": td.AccessToken})
		return
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "token expired"})
		return
	}
}

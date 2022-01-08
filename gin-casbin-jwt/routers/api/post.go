package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/models"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/repo"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/services"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/settings"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	gormadapter "github.com/casbin/gorm-adapter/v3"
)

type Post struct {
	// Required: true
	// Example: title
	Title string `json:"title" binding:"required"`
	// Required: true
	// Example: text
	Text string `json:"text" binding:"required"`
}

// swagger:route POST /api/v1/post NewPost newPostParameter
//
//	used by user to add new post
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
//	200: responseCreated
//	400: responseBadRequest
//	401: responseUnauthorized
//
// Handler function to add new post
func NewPost(c *gin.Context) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON provided"})
		return
	}
	authorization := c.Request.Header.Get("Authorization")
	content := strings.Split(authorization, " ")
	token := content[1]
	dataMap, _ := services.VerifyJWT(token, settings.AppSettings.Items.JwtAccess)

	p := models.NewPostStruct()
	p.Title = post.Title
	p.Text = post.Text

	claims, ok := dataMap.Claims.(jwt.MapClaims)
	username, ok := claims["username"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	uuid, ok := claims["access_uuid"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	redisClient := repo.GetRedisClient()
	_, err := redisClient.Get(uuid).Result()
	// user does'nt exist on redis
	if err == redis.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	
	p.Username = username
	p.NewPost()
	
	c.JSON(http.StatusCreated, gin.H{"msg": "created"})
}

// swagger:route PUT /api/v1/post/:id UpdatePost updatePostParameter
//
//	used by user to add new post
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
//	500: responseInternalServerError
//
// Handler function to update post
func UpdatePost(c *gin.Context) {
	// make gorm adapter enforcing
	db := repo.GetDB()
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(err)
	}
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON provided"})
		return
	}
	authorization := c.Request.Header.Get("Authorization")
	content := strings.Split(authorization, " ")
	token := content[1]

	dataMap, _ := services.VerifyJWT(token, settings.AppSettings.Items.JwtAccess)
	claims, ok := dataMap.Claims.(jwt.MapClaims)
	username, ok := claims["username"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	uuid, ok := claims["access_uuid"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	redisClient := repo.GetRedisClient()
	_, err = redisClient.Get(uuid).Result()
	// user does'nt exist on redis
	if err == redis.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	ID := c.Param("id")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON provided"})
		return
	}
	intID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON provided"})
		return
	}
	postModel := models.NewPostStruct()
	postModel.ID = intID
	postModel.Title = post.Title
	postModel.Text = post.Text
	newPost := postModel.GetPostById()

	ok, err = services.Enforce(username, newPost.Username, "update", adapter)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal server error"})
		c.Abort()
		return
	}
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		c.Abort()
		return
	}
	postModel.UpdatePost()

	c.JSON(http.StatusOK, gin.H{"msg": "Success"})
}

// swagger:route DELETE /api/v1/post/:id DeletePost deletePostParameter
//
//	used by user to delete post
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
//	500: responseInternalServerError
//
// Handler function to delete post
func DeletePost(c *gin.Context) {
	// make gorm adapter enforcing
	db := repo.GetDB()
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(err)
	}
	var post Post 
	if err := c.ShouldBindJSON(&post); err != nil { c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON provided"})
		return
	}
	authorization := c.Request.Header.Get("Authorization")
	content := strings.Split(authorization, " ")
	token := content[1]
	dataMap, _ := services.VerifyJWT(token, settings.AppSettings.Items.JwtAccess)
	claims, ok := dataMap.Claims.(jwt.MapClaims)
	username, ok := claims["username"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	uuid, ok := claims["access_uuid"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}
	redisClient := repo.GetRedisClient()
	_, err = redisClient.Get(uuid).Result()
	// user does'nt exist on redis
	if err == redis.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		return
	}

	ID := c.Param("id")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON provided"})
		return
	}
	intID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON provided"})
		return
	}
	postModel := models.NewPostStruct()
	postModel.ID = intID
	postModel.Title = post.Title
	postModel.Text = post.Text
	newPost := postModel.GetPostById()

	ok, err = services.Enforce(username, newPost.Username, "delete", adapter)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal server error"})
		c.Abort()
		return
	}
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		c.Abort()
		return
	}
	postModel.DeletePost()

	c.JSON(http.StatusOK, gin.H{"msg": "Suceess"})
}

// swagger:route GET /api/v1/post/:id GetPost getPostParameter
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
//	200: responseGetDataSuccess
//	400: responseBadRequest
//
// Handler function to delete post
func GetPost(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON provided"})
		return
	}
	intId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid JSON provided"})
		return
	}
	post := models.NewPostStruct()
	post.ID = intId
	postData := post.GetPostById()
	c.JSON(http.StatusOK, gin.H{"msg": "Success", "datas": postData})
}

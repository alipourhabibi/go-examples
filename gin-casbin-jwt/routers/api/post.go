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
	Title string `json:"title" binding:"required"`
	Text string `json:"text" binding:"required"`
}

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

	c.JSON(http.StatusOK, gin.H{"msg": "Updated successfully"})
}

func DeletePost(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"msg": "Deleted successfully"})
}

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
	c.JSON(http.StatusOK, gin.H{"msg": "success", "datas": postData})
}

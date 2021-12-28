package api

import (
	"net/http"
	"strings"

	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/models"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/services"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/settings"
	"github.com/gin-gonic/gin"
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
	username, _ := services.VerifyJWT(token, settings.AppSettings.Items.JwtAccess)

	p := models.NewPostStruct()
	p.Title = post.Title
	p.Text = post.Text
	p.UserName = username
	
	c.JSON(http.StatusCreated, gin.H{"msg": "created"})
}

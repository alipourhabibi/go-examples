package models

import "github.com/alipourhabibi/go-examples/gin-casbin-jwt/repo"

type Post struct {
	UserName string `json:"username"`
	Title string `json:"title"`
	Text string `json:"text"`
}

func NewPostStruct() *Post {
	return &Post{}
}

func (p *Post) NewPost() {
	newPost := repo.NewPostStruct()
	newPost.Title = p.Title
	newPost.Text = p.Text
	newPost.NewPost()
}

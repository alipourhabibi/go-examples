package models

import "github.com/alipourhabibi/go-examples/gin-casbin-jwt/repo"

type Post struct {
	ID int `json:"id"`
	Username string `json:"username"`
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
	newPost.Username = p.Username
	newPost.NewPost()
}

func (p *Post) GetPostById() repo.Post {
	post := repo.NewPostStruct()
	post.ID = p.ID
	return post.GetPostById()
}

func (p *Post) UpdatePost() {
	post := repo.NewPostStruct()
	post.ID = p.ID
	post.Username = p.Username
	post.Title = p.Title
	post.Text = p.Text
	post.UpdatePost()
}

func (p *Post) DeletePost() {
	post := repo.NewPostStruct()
	post.ID = p.ID
	post.Username = p.Username
	post.Title = p.Title
	post.Text = p.Text
	post.DeletePost()
}

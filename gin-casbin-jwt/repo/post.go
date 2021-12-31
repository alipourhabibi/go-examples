package repo

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	ID int `gorm:"id" json:"id"`
	Username string `gorm:"username" json:"username"`
	Title string `gorm:"title" json:"title"`
	Text string `gorm:"text" json:"text"`
}

func NewPostStruct() *Post {
	return &Post{}
}

func (p *Post) NewPost() {
	db := GetDB()
	db.Create(p)
}

func (p *Post) GetPostById() Post {
	post := Post{}
	db := GetDB()
	db.First(&post, "id = ? ", p.ID)
	return post
}

func (p *Post) UpdatePost() {
	db := GetDB()
	db.Model(&Post{}).Where("id = ? ", p.ID).Updates(Post{Title: p.Title, Text: p.Text})
}

func (p *Post) DeletePost() {
	db := GetDB()
	db.Model(&Post{}).Where("id = ? ", p.ID).Delete(&Post{})
}

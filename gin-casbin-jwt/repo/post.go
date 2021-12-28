package repo

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UserName string `json:"username"`
	Title string `json:"title"`
	Text string `json:"text"`
}

func NewPostStruct() *Post {
	return &Post{}
}

func (p *Post) NewPost() {
	db := GetDB()
	db.Create(p)
}

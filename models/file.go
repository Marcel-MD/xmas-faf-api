package models

type File struct {
	Base
	PostID string `json:"postId"`
	Post   Post   `json:"post" gorm:"foreignKey:PostID"`
	Name   string `json:"name"`
	Url    string `json:"url"`
	Ext    string `json:"ext"`
}

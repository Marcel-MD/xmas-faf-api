package models

type Comment struct {
	Base
	PostID string `json:"postId"`
	Post   Post   `json:"post" gorm:"foreignKey:PostID"`
	UserID string `json:"userId"`
	User   User   `json:"user" gorm:"foreignKey:UserID"`
	Text   string `json:"text"`
}

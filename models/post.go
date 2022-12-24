package models

type Post struct {
	Base
	TrainingID string    `json:"trainingId"`
	Training   Training  `json:"training" gorm:"foreignKey:TrainingID"`
	UserID     string    `json:"userId"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
	Title      string    `json:"title"`
	Text       string    `json:"text"`
	Files      []File    `json:"files" gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	Comments   []Comment `json:"comments" gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

package models

type Post struct {
	Base
	TrainingID string   `json:"trainingId"`
	Training   Training `json:"-" gorm:"foreignKey:TrainingID"`
	UserID     string   `json:"userId"`
	User       User     `json:"-" gorm:"foreignKey:UserID"`
	Title      string   `json:"title"`
	Text       string   `json:"text"`
}

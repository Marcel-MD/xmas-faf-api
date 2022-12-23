package models

type Training struct {
	Base
	Name    string `json:"name"`
	OwnerID string `json:"ownerId"`
	Users   []User `json:"users" gorm:"many2many:training_users;constraint:OnDelete:CASCADE"`
	Posts   []Post `json:"posts" gorm:"foreignKey:TrainingID;constraint:OnDelete:CASCADE"`
}

package models

import "github.com/lib/pq"

type User struct {
	Base
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email" gorm:"uniqueIndex"`
	Phone     string `json:"-"`
	Password  string `json:"-"`

	Roles pq.StringArray `json:"roles" gorm:"type:text[]"`

	Trainings []Training `json:"trainings" gorm:"many2many:training_users;constraint:OnDelete:CASCADE"`
	Comments  []Comment  `json:"comments" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	Points int `json:"points"`
}

func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}

	return false
}

const (
	UserRole  = "user"
	AdminRole = "admin"
)

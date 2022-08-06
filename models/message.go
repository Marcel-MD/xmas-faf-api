package models

const (
	Create     = "Create"
	Update     = "Update"
	Delete     = "Delete"
	RemoveUser = "RemoveUser"
	AddUser    = "AddUser"
)

type Message struct {
	Base
	RoomID   string `json:"roomId"`
	Room     Room   `json:"-" gorm:"foreignKey:RoomID"`
	UserID   string `json:"userId"`
	User     User   `json:"user" gorm:"foreignKey:UserID"`
	Text     string `json:"text"`
	Command  string `json:"command"`
	TargetID string `json:"targetId"`
}

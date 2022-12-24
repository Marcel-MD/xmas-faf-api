package dto

type CreateComment struct {
	Text string `json:"text" binding:"required,min=1,max=500"`
}

type UpdateComment struct {
	Text string `json:"text" binding:"required,min=1,max=500"`
}

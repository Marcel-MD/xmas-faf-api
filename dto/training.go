package dto

type CreateTraining struct {
	Name string `json:"name" binding:"required,min=3,max=50"`
}

type UpdateTraining struct {
	Name string `json:"name" binding:"required,min=3,max=50"`
}

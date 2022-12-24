package dto

type CreateTraining struct {
	Name     string `json:"name" binding:"required,min=3,max=50"`
	Price    int    `json:"price" binding:"required"`
	Category string `json:"category" binding:"required"`
}

type UpdateTraining struct {
	Name     string `json:"name" binding:"required,min=3,max=50"`
	Price    int    `json:"price" binding:"required"`
	Category string `json:"category" binding:"required"`
}

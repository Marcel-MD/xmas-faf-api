package dto

type CreatePost struct {
	Text  string `json:"text" binding:"required,min=1,max=500"`
	Title string `json:"title"`
}

type UpdatePost struct {
	Text  string `json:"text" binding:"required,min=1,max=500"`
	Title string `json:"title"`
}

type PostQueryParams struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

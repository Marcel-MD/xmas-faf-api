package handlers

import (
	"net/http"
	"strconv"

	"github.com/Marcel-MD/xmas-faf-api/dto"
	"github.com/Marcel-MD/xmas-faf-api/middleware"
	"github.com/Marcel-MD/xmas-faf-api/services"
	"github.com/gin-gonic/gin"
)

type postHandler struct {
	service services.IPostService
}

func routePostHandler(router *gin.RouterGroup) {
	h := &postHandler{
		service: services.GetPostService(),
	}

	r := router.Group("/posts").Use(middleware.JwtAuth())
	r.GET("/:training_id", h.find)
	r.POST("/:training_id", h.create)
	r.PUT("/:id", h.update)
	r.DELETE("/:id", h.delete)
}

func (h *postHandler) find(c *gin.Context) {
	trainingID := c.Param("training_id")
	var err error
	params := dto.PostQueryParams{}

	params.Page, err = strconv.Atoi(c.Query("page"))
	if err != nil {
		params.Page = 1
	}

	params.Size, err = strconv.Atoi(c.Query("size"))
	if err != nil {
		params.Size = 20
	}

	userID := c.GetString("user_id")

	posts, err := h.service.FindByTrainingID(trainingID, userID, params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *postHandler) create(c *gin.Context) {
	trainingID := c.Param("training_id")

	var dto dto.CreatePost
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")

	post, err := h.service.Create(trainingID, userID, dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *postHandler) update(c *gin.Context) {
	id := c.Param("id")

	var dto dto.UpdatePost
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")

	post, err := h.service.Update(id, userID, dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *postHandler) delete(c *gin.Context) {
	id := c.Param("id")

	userID := c.GetString("user_id")

	post, err := h.service.Delete(id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

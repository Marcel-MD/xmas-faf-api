package handlers

import (
	"net/http"

	"github.com/Marcel-MD/xmas-faf-api/dto"
	"github.com/Marcel-MD/xmas-faf-api/middleware"
	"github.com/Marcel-MD/xmas-faf-api/services"
	"github.com/gin-gonic/gin"
)

type trainingHandler struct {
	service     services.ITrainingService
	postService services.IPostService
}

func routeTrainingHandler(router *gin.RouterGroup) {
	h := &trainingHandler{
		service:     services.GetTrainingService(),
		postService: services.GetPostService(),
	}

	r := router.Group("/trainings")
	r.GET("/", h.findAll)
	r.GET("/:id", h.findOne)

	p := r.Use(middleware.JwtAuth())
	p.POST("/", h.create)
	p.PUT("/:id", h.update)
	p.DELETE("/:id", h.delete)
	p.POST("/:id/users/:user_id", h.addUser)
	p.DELETE("/:id/users/:user_id", h.removeUser)
}

func (h *trainingHandler) findAll(c *gin.Context) {
	trainings := h.service.FindAll()
	c.JSON(http.StatusOK, trainings)
}

func (h *trainingHandler) findOne(c *gin.Context) {
	id := c.Param("id")

	training, err := h.service.FindOne(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "training not found"})
		return
	}

	c.JSON(http.StatusOK, training)
}

func (h *trainingHandler) create(c *gin.Context) {
	userID := c.GetString("user_id")

	var dto dto.CreateTraining
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	training, err := h.service.Create(dto, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, training)
}

func (h *trainingHandler) update(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	var dto dto.UpdateTraining
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	training, err := h.service.Update(id, userID, dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, training)
}

func (h *trainingHandler) delete(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	err := h.service.Delete(id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": "training deleted"})
}

func (h *trainingHandler) addUser(c *gin.Context) {
	trainingID := c.Param("id")
	addUserID := c.Param("user_id")
	userID := c.GetString("user_id")

	err := h.service.AddUser(trainingID, addUserID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": addUserID})
}

func (h *trainingHandler) removeUser(c *gin.Context) {
	trainingID := c.Param("id")
	removeUserID := c.Param("user_id")
	userID := c.GetString("user_id")

	err := h.service.RemoveUser(trainingID, removeUserID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": removeUserID})
}

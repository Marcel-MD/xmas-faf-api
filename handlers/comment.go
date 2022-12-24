package handlers

import (
	"net/http"

	"github.com/Marcel-MD/xmas-faf-api/dto"
	"github.com/Marcel-MD/xmas-faf-api/middleware"
	"github.com/Marcel-MD/xmas-faf-api/services"
	"github.com/gin-gonic/gin"
)

type commentHandler struct {
	service services.ICommentService
}

func routeCommentHandler(router *gin.RouterGroup) {
	h := &commentHandler{
		service: services.GetCommentService(),
	}

	r := router.Group("/comments").Use(middleware.JwtAuth())
	r.GET("/:post_id", h.find)
	r.POST("/:post_id", h.create)
	r.PUT("/:id", h.update)
	r.DELETE("/:id", h.delete)
}

func (h *commentHandler) find(c *gin.Context) {
	postId := c.Param("post_id")

	userID := c.GetString("user_id")

	comments, err := h.service.FindByPostID(postId, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *commentHandler) create(c *gin.Context) {
	postID := c.Param("post_id")

	var dto dto.CreateComment
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")

	comment, err := h.service.Create(postID, userID, dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (h *commentHandler) update(c *gin.Context) {
	id := c.Param("id")

	var dto dto.UpdateComment
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")

	comment, err := h.service.Update(id, userID, dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (h *commentHandler) delete(c *gin.Context) {
	id := c.Param("id")

	userID := c.GetString("user_id")

	comment, err := h.service.Delete(id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

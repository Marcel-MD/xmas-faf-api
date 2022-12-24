package handlers

import (
	"io"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Marcel-MD/xmas-faf-api/middleware"
	"github.com/Marcel-MD/xmas-faf-api/services"
	"github.com/gin-gonic/gin"
)

type fileHandler struct {
	service     services.IFileService
	blobService services.IBlobService
}

func routeFileHandler(router *gin.RouterGroup) {
	h := &fileHandler{
		service:     services.GetFileService(),
		blobService: services.GetBlobService(),
	}

	r := router.Group("/files")
	r.GET("/:post_id", h.find)
	r.GET("/file/:file_name", h.findFile)

	a := r.Use(middleware.JwtAuth())
	a.POST("/:post_id", h.create)
	a.DELETE("/:id", h.delete)
}

func (h *fileHandler) find(c *gin.Context) {
	postID := c.Param("post_id")

	files := h.service.FindByPostID(postID)

	c.JSON(200, files)
}

func (h *fileHandler) create(c *gin.Context) {
	postID := c.Param("post_id")

	form, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	file, err := form.Open()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	h.service.Create(postID, form.Filename, data)
}

func (h *fileHandler) delete(c *gin.Context) {
	id := c.Param("id")

	h.service.Delete(id)

	c.JSON(200, gin.H{"message": "File deleted"})
}

func (h *fileHandler) findFile(c *gin.Context) {
	fileName := c.Param("file_name")

	downloadResponse, err := h.blobService.Get(fileName)

	defer downloadResponse.Body(azblob.RetryReaderOptions{}).Close()

	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)

	_, err = io.Copy(c.Writer, downloadResponse.Body(azblob.RetryReaderOptions{}))

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

}

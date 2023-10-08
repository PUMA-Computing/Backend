package files

import (
	"Backend/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type Handler struct {
	PermissionService *services.PermissionService
	FilesService      *services.FilesService
}

func NewFilesHandler(filesServices *services.FilesService, permissionService *services.PermissionService) *Handler {
	return &Handler{FilesService: filesServices, PermissionService: permissionService}
}

func (h *Handler) UploadFile(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Failed to upload file"})
	}
	fileName := h.FilesService.ConvertFileName(file.Filename)
	saveFile := h.FilesService.ConvertToPublicDirectory(fileName)
	if err := c.SaveUploadedFile(file, saveFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"type": "URL",
		"url":  "http://localhost:" + os.Getenv("SERVER_PORT") + "/public/" + fileName,
	}})
	return
}

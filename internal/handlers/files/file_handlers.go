package files

import (
	"Backend/internal/services"
	"Backend/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handlers struct {
	FilesService      *services.FilesService
	PermissionService *services.PermissionService
}

func NewFilesHandlers(filesService *services.FilesService, permissionService *services.PermissionService) *Handlers {
	return &Handlers{
		FilesService:      filesService,
		PermissionService: permissionService,
	}
}

func (h *Handlers) UploadFile(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "files:upload")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": []string{"Permission Denied"}})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to upload file"})
	}

	if !h.FilesService.IsImageFile(file) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "File is not an image"})
		return
	}

	if !h.FilesService.IsFileSizeValid(file) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Maximum File Size is 20MB"})
		return
	}

	if !h.FilesService.IsFileExtensionValid(file) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "File extension is not allowed"})
		return
	}

	fileName := utils.GenerateUniqueFileName(file.Filename)
	err = h.FilesService.UploadFileToBucket(file, fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Cannot upload file to buckets"})
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"success": "true",
		"message": "File uploaded successfully",
		"data": gin.H{
			"type":      "url",
			"file_name": fileName,
		},
	}})
	return
}

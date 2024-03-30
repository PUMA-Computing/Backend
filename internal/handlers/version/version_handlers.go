package version

import (
	"Backend/internal/models"
	"Backend/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handlers struct {
	VersionService *services.VersionService
}

func NewVersionHandlers(versionService *services.VersionService) *Handlers {
	return &Handlers{VersionService: versionService}
}

func (h *Handlers) GetVersion(c *gin.Context) {
	version, err := h.VersionService.GetVersion()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"version": version})
}

func (h *Handlers) UpdateVersion(c *gin.Context) {
	var version models.Version
	if err := c.ShouldBindJSON(&version); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.VersionService.UpdateVersion(&version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"version": version})
}

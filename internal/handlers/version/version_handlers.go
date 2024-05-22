package version

import (
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

// GetVersion retrieves the current version from the database
func (h *Handlers) GetVersion(c *gin.Context) {
	version, err := h.VersionService.GetVersion()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"latest_version": version})
}

// GetChangelog retrieves the changelog from the database
func (h *Handlers) GetChangelog(c *gin.Context) {
	changelog, err := h.VersionService.GetChangelog()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"changelog": changelog})
}

package news

import (
	"Backend/internal/models"
	"Backend/internal/services"
	"Backend/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handler struct {
	NewsService       *services.NewsService
	PermissionService *services.PermissionService
}

func NewNewsHandler(newsService *services.NewsService, permissionService *services.PermissionService) *Handler {
	return &Handler{
		NewsService:       newsService,
		PermissionService: permissionService,
	}
}

func (h *Handler) CreateNews(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "news:create")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"Permission Denied"}})
		return
	}

	var newNews models.News
	if err := c.BindJSON(&newNews); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	if err := h.NewsService.CreateNews(&newNews); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": newNews,
		"meta": gin.H{"message": "News Created Successfully"},
	})
}

func (h *Handler) GetNewsByID(c *gin.Context) {
	newsIDStr := c.Param("newsID")
	newsID, err := strconv.Atoi(newsIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Invalid News ID"}})
		return
	}

	news, err := h.NewsService.GetNewsByID(newsID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": []string{"News not found"}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": news})
}

func (h *Handler) EditNews(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "news:edit")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"Permission Denied"}})
		return
	}

	newsIDStr := c.Param("newsID")
	newsID, err := strconv.Atoi(newsIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Invalid News ID"}})
		return
	}
	var updatedNews models.News
	if err := c.BindJSON(&updatedNews); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	if err := h.NewsService.EditNews(newsID, &updatedNews); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
	}
	c.JSON(http.StatusOK, gin.H{"data": updatedNews})
}

func (h *Handler) DeleteNews(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "news:delete")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"errors": []string{"Permission Denied"}})
		return
	}

	newsIDStr := c.Param("newsID")
	newsID, err := strconv.Atoi(newsIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Invalid News ID"}})
		return
	}

	if err := h.NewsService.DeleteNews(newsID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"meta": gin.H{"message": "News Deleted Successfully"}})
}

func (h *Handler) ListNews(c *gin.Context) {
	news, err := h.NewsService.ListNews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": news})
}

func (h *Handler) LikeNews(c *gin.Context) {
	newsIDStr := c.Param("newsID")
	newsID, err := strconv.Atoi(newsIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Invalid News ID"}})
		return
	}

	if err := h.NewsService.LikeNews(newsID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"meta": gin.H{"message": "News Liked Successfully"}})
}

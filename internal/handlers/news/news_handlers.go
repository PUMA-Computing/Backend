package news

import (
	"Backend/internal/models"
	"Backend/internal/services"
	"Backend/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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
	token, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}
	userID, err := utils.GetUserIDFromToken(token, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "news:create")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": []string{"Permission Denied"}})
		return
	}

	var newNews models.News
	if err := c.BindJSON(&newNews); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	newNews.UserID = userID
	newNews.CreatedAt = time.Time{}
	newNews.UpdatedAt = time.Time{}

	// slug from title
	newNews.Slug = utils.GenerateFriendlyURL(newNews.Title)

	if err := h.NewsService.CreateNews(&newNews); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "News Created Successfully",
		"data":    newNews,
	})
}

func (h *Handler) GetNewsByID(c *gin.Context) {
	newsIDStr := c.Param("newsID")
	newsID, err := strconv.Atoi(newsIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid News ID"}})
		return
	}

	news, err := h.NewsService.GetNewsByID(newsID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": []string{"News not found"}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "News Retrieved Successfully",
		"data":    news,
	})
}

func (h *Handler) EditNews(c *gin.Context) {
	token, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}
	userID, err := utils.GetUserIDFromToken(token, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "news:edit")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": []string{"Permission Denied"}})
		return
	}

	newsIDStr := c.Param("newsID")
	newsID, err := strconv.Atoi(newsIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid News ID"}})
		return
	}

	var updatedNews models.News
	if err := c.BindJSON(&updatedNews); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	updatedNews.UpdatedAt = time.Time{}

	updatedAttributes := make(map[string]interface{})
	if updatedNews.Title != "" {
		updatedAttributes["title"] = updatedNews.Title
	}

	if updatedNews.Content != "" {
		updatedAttributes["content"] = updatedNews.Content
	}

	if updatedNews.PublishDate.IsZero() {
		updatedAttributes["publish_date"] = updatedNews.PublishDate
	}

	if updatedNews.Likes != 0 {
		updatedAttributes["likes"] = updatedNews.Likes
	}

	if err := h.NewsService.EditNews(newsID, &updatedNews); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "News Updated Successfully",
		"data":    updatedAttributes,
	})
}

func (h *Handler) DeleteNews(c *gin.Context) {
	token, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}
	userID, err := utils.GetUserIDFromToken(token, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	hasPermission, err := h.PermissionService.CheckPermission(context.Background(), userID, "news:delete")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": []string{"Permission Denied"}})
		return
	}

	newsIDStr := c.Param("newsID")
	newsID, err := strconv.Atoi(newsIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid News ID"}})
		return
	}

	if err := h.NewsService.DeleteNews(newsID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "News Deleted Successfully",
	})
}

func (h *Handler) ListNews(c *gin.Context) {

	news, err := h.NewsService.ListNews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	log.Println("List Events End")

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"totalResults": len(news),
		"data":         news,
	})
}

func (h *Handler) LikeNews(c *gin.Context) {
	token, err := utils.ExtractTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}
	userID, err := utils.GetUserIDFromToken(token, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	newsIDStr := c.Param("newsID")
	newsID, err := strconv.Atoi(newsIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid News ID"}})
		return
	}

	if err := h.NewsService.LikeNews(userID, newsID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "News Liked Successfully",
	})
}

//func (h *Handler) UnlikeNews(c *gin.Context) {
//	token, err := utils.ExtractTokenFromHeader(c)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
//		return
//	}
//	userID, err := utils.GetUserIDFromToken(token, os.Getenv("JWT_SECRET_KEY"))
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
//		return
//	}
//
//	newsIDStr := c.Param("newsID")
//	newsID, err := strconv.Atoi(newsIDStr)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid News ID"}})
//		return
//	}
//
//	if err := h.NewsService.UnlikeNews(userID, newsID); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"success": true,
//		"message": "News Unliked Successfully",
//	})
//}

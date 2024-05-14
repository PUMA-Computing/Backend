package news

import (
	"Backend/internal/handlers/auth"
	"Backend/internal/models"
	"Backend/internal/services"
	"Backend/pkg/utils"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Handler struct {
	NewsService       *services.NewsService
	PermissionService *services.PermissionService
	AWSService        *services.S3Service
	R2Service         *services.S3Service
}

func NewNewsHandler(newsService *services.NewsService, permissionService *services.PermissionService, AWSService *services.S3Service, R2Service *services.S3Service) *Handler {
	return &Handler{
		NewsService:       newsService,
		PermissionService: permissionService,
		AWSService:        AWSService,
		R2Service:         R2Service,
	}
}

func (h *Handler) CreateNews(c *gin.Context) {
	userID, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "news:create")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	data := c.Request.FormValue("data")
	var newNews models.News
	if err := json.Unmarshal([]byte(data), &newNews); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	newNews.UserID = userID

	if newNews.Title == "" {
		newNews.Slug = utils.GenerateFriendlyURL(newNews.Title)
	}

	if newNews.PublishDate.IsZero() {
		newNews.PublishDate = time.Now()
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	optimizedImage, err := utils.OptimizeImage(file, 2800, 1080)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	optimizedImageBytes, err := io.ReadAll(optimizedImage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	// Choose storage service to upload image to (AWS or R2)
	upload := utils.ChooseStorageService()

	if upload == utils.R2Service {
		err = h.R2Service.UploadFileToR2(context.Background(), "news", newNews.Slug, optimizedImageBytes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}

		newNews.Thumbnail, _ = h.R2Service.GetFileR2("news", newNews.Slug)
	} else {
		err = h.AWSService.UploadFileToAWS(context.Background(), "news", newNews.Slug, optimizedImageBytes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}

		newNews.Thumbnail, _ = h.AWSService.GetFileAWS("news", newNews.Slug)
	}

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
	_, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "news:edit")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	newsIDStr := c.Param("newsID")
	newsID, err := strconv.Atoi(newsIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{"Invalid News ID"}})
		return
	}

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	data := c.Request.FormValue("data")
	var updatedNews models.News
	if err := json.Unmarshal([]byte(data), &updatedNews); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	optimizedImage, err := utils.OptimizeImage(file, 2800, 1080)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	optimizedImageBytes, err := io.ReadAll(optimizedImage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	// Choose storage service to upload image to (AWS or R2)
	upload := utils.ChooseStorageService()

	if upload == utils.R2Service {
		err = h.R2Service.UploadFileToR2(context.Background(), "news", updatedNews.Slug, optimizedImageBytes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}

		updatedNews.Thumbnail, _ = h.R2Service.GetFileR2("news", updatedNews.Slug)
	} else {
		err = h.AWSService.UploadFileToAWS(context.Background(), "news", updatedNews.Slug, optimizedImageBytes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}

		updatedNews.Thumbnail, _ = h.AWSService.GetFileAWS("news", updatedNews.Slug)
	}

	existingNews, err := h.NewsService.GetNewsByID(newsID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": []string{"News not found"}})
		return
	}

	if updatedNews.Title == "" {
		updatedNews.Slug = utils.GenerateFriendlyURL(updatedNews.Title)
	} else {
		updatedNews.Slug = existingNews.Slug
	}

	utils.ReflectiveUpdate(existingNews, &updatedNews)

	if err := h.NewsService.EditNews(newsID, &updatedNews); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "News Updated Successfully",
		"data":    existingNews,
	})
}

func (h *Handler) DeleteNews(c *gin.Context) {
	_, err := (&auth.Handlers{}).ExtractUserIDAndCheckPermission(c, "news:delete")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

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

	// Check image exists on AWS or R2
	exists, _ := h.AWSService.FileExists(context.Background(), "news", news.Slug)
	if exists {
		if err := h.AWSService.DeleteFile(context.Background(), "news", news.Slug); err != nil {
			log.Println("Error deleting file from AWS")
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
			return
		}
	} else {
		exists, _ := h.R2Service.FileExists(context.Background(), "news", news.Slug)
		if exists {
			if err := h.R2Service.DeleteFile(context.Background(), "news", news.Slug); err != nil {
				log.Println("Error deleting file from R2")
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
				return
			}
		}
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
	queryParams := make(map[string]string)
	queryParams["organization_id"] = c.Query("organization_id")
	queryParams["page"] = c.Query("page")

	news, totalPages, err := h.NewsService.ListNews(queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"data":         news,
		"totalResults": len(news),
		"totalPages":   totalPages,
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

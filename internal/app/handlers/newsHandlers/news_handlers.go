package newsHandlers

import (
	"Backend/internal/app/domain/news"
	"Backend/internal/app/interfaces/service/newsService"
	"Backend/internal/utils/validation"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type NewsHandlers struct {
	newsService newsService.NewsService
}

func NewNewsHandlers(newsService newsService.NewsService) *NewsHandlers {
	return &NewsHandlers{newsService: newsService}
}

func (h *NewsHandlers) GetNews() fiber.Handler {
	return func(c *fiber.Ctx) error {
		newsList, err := h.newsService.GetNews()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error retrieving News", "error": err.Error()})
		}
		return c.JSON(newsList)
	}
}

func (h *NewsHandlers) GetNewsByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		newsID, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid news ID"})
		}

		newsItem, err := h.newsService.GetNewsByID(int64(newsID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error retrieving news by ID"})
		}
		return c.JSON(newsItem)
	}
}

func (h *NewsHandlers) CreateNews() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var newsItem news.News
		if err := c.BodyParser(&newsItem); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error parsing news"})
		}

		if err := validation.ValidateNews(&newsItem); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error validating news"})
		}

		if err := h.newsService.CreateNews(&newsItem); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error creating news"})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "News created successfully"})
	}
}

func (h *NewsHandlers) EditNews() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var newsItem news.News
		if err := c.BodyParser(&newsItem); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error parsing news"})
		}

		if err := validation.ValidateNews(&newsItem); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error validating news"})
		}

		if err := h.newsService.UpdateNews(&newsItem); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error editing news"})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "News edited successfully"})
	}
}

func (h *NewsHandlers) DeleteNews() fiber.Handler {
	return func(c *fiber.Ctx) error {
		newsID := c.Params("id")
		newsIDInt, err := strconv.ParseInt(newsID, 10, 64)

		err = h.newsService.DeleteNews(newsIDInt)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error deleting news"})
		}

		return c.JSON(fiber.Map{"message": "News deleted Successfully"})
	}
}

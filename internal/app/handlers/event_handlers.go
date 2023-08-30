package handlers

import (
	"Backend/internal/app/service"
	"Backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type EventHandlers struct {
	eventService service.EventService
}

func NewEventHandlers(eventService service.EventService) *EventHandlers {
	return &EventHandlers{eventService: eventService}
}

func (h *EventHandlers) GetEventByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		events, err := h.eventService.GetEvent()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error retrieving event"})
		}
		return c.JSON(events)
	}
}

func (h *EventHandlers) RegisterUserForEvent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		eventID := c.Params("eventID")
		userID := utils.GetUserIDFromContext(c)
		if err := h.eventService.RegisterUserForEvent(strconv.Itoa(int(userID)), eventID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error registering user for event"})
		}
		return c.JSON(fiber.Map{"message": "User registered for event"})
	}
}

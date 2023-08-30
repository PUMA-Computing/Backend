package handlers

import (
	"Backend/internal/app/domain"
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

func (h *EventHandlers) CreateEvent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var event domain.Event
		if err := c.BodyParser(&event); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error Parsing event"})
		}

		err := h.eventService.CreateEvent(&event)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error creating event"})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Event created Successfully"})
	}
}

func (h *EventHandlers) EditEvent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		eventID := c.Params("id")
		var UpdatedEvent domain.Event
		if err := c.BodyParser(&UpdatedEvent); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error Parsing Updated event"})
		}

		err := h.eventService.UpdateEvent(eventID, &UpdatedEvent)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error updating event"})
		}

		return c.JSON(fiber.Map{"message": "Event updated Successfully"})
	}
}

func (h *EventHandlers) DeleteEvent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		eventID := c.Params("id")

		err := h.eventService.DeleteEvent(eventID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error deleting event"})
		}

		return c.JSON(fiber.Map{"message": "Event deleted Successfully"})
	}
}

func (h *EventHandlers) GetEvent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		events, err := h.eventService.GetEvent()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error retrieving events"})
		}
		return c.JSON(events)
	}
}

func (h *EventHandlers) GetEventUsers() fiber.Handler {
	return func(c *fiber.Ctx) error {
		eventID := c.Params("id")

		users, err := h.eventService.GetEventUsers(eventID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error retrieving event users"})
		}

		return c.JSON(users)
	}
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

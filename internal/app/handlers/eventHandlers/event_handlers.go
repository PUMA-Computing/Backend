package eventHandlers

import (
	"Backend/internal/app/domain/event"
	"Backend/internal/app/interfaces/service/eventService"
	"Backend/internal/utils/getUserContext"
	"Backend/internal/utils/validation"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type EventHandlers struct {
	eventService eventService.EventService
}

func NewEventHandlers(eventService eventService.EventService) *EventHandlers {
	return &EventHandlers{eventService: eventService}
}

func (h *EventHandlers) CreateEvent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var event event.Event
		if err := c.BodyParser(&event); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error Parsing event"})
		}

		/*
			Validate eventService data
		*/
		if err := validation.ValidateEvent(&event); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error validating event"})
		}

		err := h.eventService.CreateEvent(&event)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error creating event", "error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Event created Successfully"})
	}
}

func (h *EventHandlers) EditEvent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		eventID := c.Params("id")
		eventIDInt, err := strconv.ParseInt(eventID, 10, 64)
		var UpdatedEvent event.Event
		if err := c.BodyParser(&UpdatedEvent); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error Parsing Updated event"})
		}

		err = h.eventService.UpdateEvent(eventIDInt, &UpdatedEvent)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error updating event"})
		}

		return c.JSON(fiber.Map{"message": "Event updated Successfully", "error": err.Error()})
	}
}

func (h *EventHandlers) DeleteEvent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		eventID := c.Params("id")
		eventIDInt, err := strconv.ParseInt(eventID, 10, 64)

		err = h.eventService.DeleteEvent(eventIDInt)
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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error retrieving events", "errors": err.Error()})
		}
		return c.JSON(events)
	}
}

func (h *EventHandlers) GetEventUsers() fiber.Handler {
	return func(c *fiber.Ctx) error {
		eventID := c.Params("id")
		eventIDInt, err := strconv.ParseInt(eventID, 10, 64)

		users, err := h.eventService.GetEventUsers(eventIDInt)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error retrieving event users"})
		}

		return c.JSON(users)
	}
}

func (h *EventHandlers) GetEventByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		eventID := c.Params("id")
		eventIDInt, err := strconv.ParseInt(eventID, 10, 64)

		event, err := h.eventService.GetEventByID(eventIDInt)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error retrieving event"})
		}

		return c.JSON(event)
	}
}

func (h *EventHandlers) RegisterUserForEvent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		eventID := c.Params("eventID")
		userID, _ := getUserContext.GetUserIDFromContext(c)
		eventIDInt, err := strconv.ParseInt(eventID, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error parsing eventID"})
		}
		if err := h.eventService.RegisterUserForEvent(userID, eventIDInt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error registering user for event"})
		}
		return c.JSON(fiber.Map{"message": "User registered for event"})
	}
}

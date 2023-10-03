package handlers

import (
	"Backend/internal/services"
)

type Handlers struct {
	UserService       *services.UserService
	EventService      *services.EventService
	NewsService       *services.NewsService
	RoleService       *services.RoleService
	PermissionService *services.PermissionService
}

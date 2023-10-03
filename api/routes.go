package api

import (
	"Backend/api/handlers/event"
	"Backend/api/handlers/news"
	"Backend/api/handlers/permission"
	"Backend/api/handlers/role"
	"Backend/api/handlers/user"
	"Backend/internal/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	userService := services.NewUserService()
	eventService := services.NewEventService()
	newsService := services.NewNewsService()
	roleService := services.NewRoleService()
	permissionService := services.NewPermissionService()

	userHandlers := user.NewUserHandlers(userService)
	eventHandlers := event.NewEventHandlers(eventService)
	newsHandlers := news.NewNewsHandler(newsService)
	roleHandlers := role.NewRoleHandler(roleService)
	permissionHandlers := permission.NewPermissionHandler(permissionService)

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", userHandlers.RegisterUser)
		authRoutes.POST("/login", userHandlers.Login)
	}

	userRoutes := r.Group("/users")
	{
		userRoutes.GET("/:userID", userHandlers.GetUserByID)
		userRoutes.PUT("/:userID/edit", userHandlers.EditUser)
		userRoutes.DELETE("/:userID/delete", userHandlers.DeleteUser)
		userRoutes.GET("/", userHandlers.ListUsers)
	}

	eventRoutes := r.Group("/events")
	{
		eventRoutes.POST("/create", eventHandlers.CreateEvent)
		eventRoutes.GET("/:eventID", eventHandlers.GetEventByID)
		eventRoutes.PUT("/:eventID/edit", eventHandlers.EditEvent)
		eventRoutes.DELETE("/:eventID/delete", eventHandlers.DeleteEvent)
		eventRoutes.GET("/", eventHandlers.ListEvents)
		eventRoutes.POST("/:eventID/register", eventHandlers.RegisterForEvent)
		eventRoutes.GET("/:eventID/registered-user", eventHandlers.ListRegisteredUsers)
	}

	newsRoutes := r.Group("/news")
	{
		newsRoutes.POST("/create", newsHandlers.CreateNews)
		newsRoutes.GET("/:newsID", newsHandlers.GetNewsByID)
		newsRoutes.PUT("/:newsID/edit", newsHandlers.EditNews)
		newsRoutes.DELETE("/:newsID/delete", newsHandlers.DeleteNews)
		newsRoutes.GET("/", newsHandlers.ListNews)
		newsRoutes.POST("/:newsID/like", newsHandlers.LikeNews)
	}

	roleRoutes := r.Group("/roles")
	{
		roleRoutes.POST("/create", roleHandlers.CreateRole)
		roleRoutes.GET("/:roleID", roleHandlers.GetRoleByID)
		roleRoutes.PUT("/:roleID/edit", roleHandlers.EditRole)
		roleRoutes.DELETE("/:roleID/delete", roleHandlers.DeleteRole)
		roleRoutes.GET("/", roleHandlers.ListRoles)
		roleRoutes.POST("/:roleID/assign-permission", roleHandlers.AssignRoleToUser)
	}

	permissionRoutes := r.Group("/permissions")
	{
		permissionRoutes.GET("/", permissionHandlers.ListPermissions)
		permissionRoutes.POST("/assign/:roleID", permissionHandlers.AssingPermissionToRole)

	}
	return r
}

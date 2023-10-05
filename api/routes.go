package api

import (
	"Backend/api/handlers/event"
	"Backend/api/handlers/news"
	"Backend/api/handlers/permission"
	"Backend/api/handlers/role"
	"Backend/api/handlers/user"
	"Backend/api/middleware"
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

	userHandlers := user.NewUserHandlers(userService, permissionService)
	eventHandlers := event.NewEventHandlers(eventService, permissionService)
	newsHandlers := news.NewNewsHandler(newsService, permissionService)
	roleHandlers := role.NewRoleHandler(roleService, userService, permissionService)
	permissionHandlers := permission.NewPermissionHandler(permissionService)

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", userHandlers.RegisterUser)
		authRoutes.POST("/login", userHandlers.Login)
		authRoutes.POST("/logout", userHandlers.Logout)
		authRoutes.POST("/refresh-token", middleware.TokenMiddleware(), userHandlers.RefreshToken)
	}

	userRoutes := r.Group("/user")
	{
		userRoutes.Use(middleware.TokenMiddleware())
		userRoutes.GET("/:userID", userHandlers.GetUserByID)
		userRoutes.PUT("/edit", userHandlers.EditUser)
		userRoutes.DELETE("/delete", userHandlers.DeleteUser)
		userRoutes.GET("/list", userHandlers.ListUsers)
	}

	eventRoutes := r.Group("/event")
	{
		eventRoutes.GET("/:eventID", eventHandlers.GetEventByID)
		eventRoutes.GET("/", eventHandlers.ListEvents)
		eventRoutes.Use(middleware.TokenMiddleware())
		eventRoutes.POST("/create", eventHandlers.CreateEvent)
		eventRoutes.PUT("/:eventID/edit", eventHandlers.EditEvent)
		eventRoutes.DELETE("/:eventID/delete", eventHandlers.DeleteEvent)
		eventRoutes.POST("/:eventID/register", eventHandlers.RegisterForEvent)
		eventRoutes.GET("/:eventID/registered-users", eventHandlers.ListRegisteredUsers)
	}

	newsRoutes := r.Group("/news")
	{
		newsRoutes.GET("/:newsID", newsHandlers.GetNewsByID)
		newsRoutes.GET("/", newsHandlers.ListNews)
		newsRoutes.Use(middleware.TokenMiddleware())
		newsRoutes.POST("/create", newsHandlers.CreateNews)
		newsRoutes.PUT("/:newsID/edit", newsHandlers.EditNews)
		newsRoutes.DELETE("/:newsID/delete", newsHandlers.DeleteNews)
		newsRoutes.POST("/:newsID/like", newsHandlers.LikeNews)
	}

	roleRoutes := r.Group("/roles")
	{
		roleRoutes.Use(middleware.TokenMiddleware())
		roleRoutes.POST("/create", roleHandlers.CreateRole)
		roleRoutes.PUT("/:roleID/edit", roleHandlers.EditRole)
		roleRoutes.DELETE("/:roleID/delete", roleHandlers.DeleteRole)
		roleRoutes.GET("/", roleHandlers.ListRoles)
		roleRoutes.POST("/:roleID/assign/:userID", roleHandlers.AssignRoleToUser)
	}
	roleRoutes.GET("/:roleID", roleHandlers.GetRoleByID)

	permissionRoutes := r.Group("/permissions")
	{
		permissionRoutes.Use(middleware.TokenMiddleware())
		permissionRoutes.GET("/list", permissionHandlers.ListPermissions)
		permissionRoutes.POST("/assign/:roleID", permissionHandlers.AssignPermissionToRole)

	}
	return r
}

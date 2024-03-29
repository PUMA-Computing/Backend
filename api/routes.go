package api

import (
	"Backend/internal/handlers/aspirations"
	"Backend/internal/handlers/auth"
	"Backend/internal/handlers/event"
	"Backend/internal/handlers/files"
	"Backend/internal/handlers/news"
	"Backend/internal/handlers/permission"
	"Backend/internal/handlers/role"
	"Backend/internal/handlers/user"
	"Backend/internal/middleware"
	"Backend/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://computing.president.ac.id", "https://staging.computing.president.ac.id", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Static("/public", "./public")

	authService := services.NewAuthService()
	userService := services.NewUserService()
	eventService := services.NewEventService()
	newsService := services.NewNewsService()
	roleService := services.NewRoleService()
	permissionService := services.NewPermissionService()
	filesService := services.NewFilesService()
	aspirationsService := services.NewAspirationService()
	s3Service, _ := services.NewS3Service()

	authHandlers := auth.NewAuthHandlers(authService, permissionService)
	userHandlers := user.NewUserHandlers(userService, permissionService)
	eventHandlers := event.NewEventHandlers(eventService, permissionService, s3Service)
	newsHandlers := news.NewNewsHandler(newsService, permissionService)
	roleHandlers := role.NewRoleHandler(roleService, userService, permissionService)
	permissionHandlers := permission.NewPermissionHandler(permissionService)
	filesHandlers := files.NewFilesHandlers(filesService, permissionService)
	aspirationHandlers := aspirations.NewAspirationHandlers(aspirationsService, permissionService)

	api := r.Group("/api/v1")

	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", authHandlers.RegisterUser)
		authRoutes.POST("/login", authHandlers.Login)
		authRoutes.POST("/logout", authHandlers.Logout)
		authRoutes.POST("/refresh-token", middleware.TokenMiddleware(), authHandlers.RefreshToken)
	}

	userRoutes := api.Group("/user")
	{
		userRoutes.Use(middleware.TokenMiddleware())
		userRoutes.GET("/:userID", userHandlers.GetUserByID)
		userRoutes.PUT("/edit", userHandlers.EditUser)
		userRoutes.DELETE("/delete", userHandlers.DeleteUser)
		userRoutes.GET("/list", userHandlers.ListUsers)

		// ListEventsRegisteredByUser
		userRoutes.GET("/registered-events", eventHandlers.ListEventsRegisteredByUser)
	}

	eventRoutes := api.Group("/event")
	{
		eventRoutes.GET("/:eventID", eventHandlers.GetEventBySlug)
		eventRoutes.GET("/", eventHandlers.ListEvents)
		eventRoutes.Use(middleware.TokenMiddleware())
		eventRoutes.POST("/create", eventHandlers.CreateEvent)
		eventRoutes.PATCH("/:eventID/edit", eventHandlers.EditEvent)
		eventRoutes.DELETE("/:eventID/delete", eventHandlers.DeleteEvent)
		eventRoutes.POST("/:eventID/register", eventHandlers.RegisterForEvent)
		eventRoutes.GET("/:eventID/registered-users", eventHandlers.ListRegisteredUsers)
	}

	newsRoutes := api.Group("/news")
	{
		newsRoutes.GET("/", newsHandlers.ListNews)
		newsRoutes.GET("/:newsID", newsHandlers.GetNewsByID)
		newsRoutes.Use(middleware.TokenMiddleware())
		newsRoutes.POST("/create", newsHandlers.CreateNews)
		newsRoutes.PUT("/:newsID/edit", newsHandlers.EditNews)
		newsRoutes.DELETE("/:newsID/delete", newsHandlers.DeleteNews)
		newsRoutes.POST("/:newsID/like", newsHandlers.LikeNews)
	}

	roleRoutes := api.Group("/roles")
	{
		roleRoutes.Use(middleware.TokenMiddleware())
		roleRoutes.GET("/", roleHandlers.ListRoles)
		roleRoutes.POST("/create", roleHandlers.CreateRole)
		roleRoutes.GET("/:roleID", roleHandlers.GetRoleByID)
		roleRoutes.PUT("/:roleID/edit", roleHandlers.EditRole)
		roleRoutes.DELETE("/:roleID/delete", roleHandlers.DeleteRole)
		roleRoutes.POST("/:roleID/assign/:userID", roleHandlers.AssignRoleToUser)
	}
	permissionRoutes := api.Group("/permissions")
	{
		permissionRoutes.Use(middleware.TokenMiddleware())
		permissionRoutes.GET("/list", permissionHandlers.ListPermissions)
		permissionRoutes.POST("/assign/:roleID", permissionHandlers.AssignPermissionToRole)

	}

	filesRoutes := api.Group("/files")
	{
		//filesRoutes.PUT("/", filesHandlers.UploadFile)
		filesRoutes.PUT("/upload/profile-picture", filesHandlers.UploadFile)
	}

	aspirationRoutes := api.Group("/aspirations")
	{
		aspirationRoutes.GET("/", aspirationHandlers.GetAspirations)
		aspirationRoutes.Use(middleware.TokenMiddleware())
		aspirationRoutes.POST("/create", aspirationHandlers.CreateAspiration)
		aspirationRoutes.PATCH("/:id/close", aspirationHandlers.CloseAspiration)
		aspirationRoutes.DELETE("/:id/delete", aspirationHandlers.DeleteAspiration)
		aspirationRoutes.POST("/:id/upvote", aspirationHandlers.UpvoteAspiration)
		aspirationRoutes.GET("/:id/get_upvotes", aspirationHandlers.GetUpvotesByAspirationID)
		aspirationRoutes.POST("/:id/admin_reply", aspirationHandlers.AddAdminReply)
	}
	return r
}

package routes

import (
	"github.com/Sparker0i/Cactro-Backend-09-Mar-25/config"
	"github.com/Sparker0i/Cactro-Backend-09-Mar-25/internal/api/handlers"
	"github.com/Sparker0i/Cactro-Backend-09-Mar-25/internal/api/middleware"
	"github.com/Sparker0i/Cactro-Backend-09-Mar-25/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the API routes
func SetupRouter(config *config.Config) *gin.Engine {
	// Set Gin mode
	gin.SetMode(config.Server.GinMode)

	// Create router with default middleware
	router := gin.New()

	// Add custom middleware
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.CORS.AllowOrigins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(corsConfig))

	// Create services
	var githubService services.GitHubServiceInterface = services.NewGitHubService(config)

	// Create handlers
	githubHandler := handlers.NewGitHubHandler(githubService)

	// Add health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// GitHub routes
	githubGroup := router.Group("/github")
	{
		githubGroup.GET("", githubHandler.GetUserProfile)
		githubGroup.GET("/:repo", githubHandler.GetRepository)
		githubGroup.POST("/:repo/issues", githubHandler.CreateIssue)
	}

	return router
}

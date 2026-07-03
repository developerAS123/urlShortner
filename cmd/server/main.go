package main

import (
	"log"
	"os"

	"github.com/ankitsingh/urlshortener/internal/handlers"
	"github.com/ankitsingh/urlshortener/internal/middleware"
	"github.com/ankitsingh/urlshortener/internal/repository"
	"github.com/ankitsingh/urlshortener/internal/services"
	"github.com/ankitsingh/urlshortener/internal/workers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (for local development)
	_ = godotenv.Load()

	// Initialize Database, Redis, and Services
	repository.InitDB()
	repository.InitRedis()
	
	services.InitGeoIP()
	defer services.CloseGeoIP()
	
	middleware.InitRateLimiter()

	r := gin.Default()

	// CORS Setup
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:5173", // Local frontend
	}
	if vercelURL := os.Getenv("FRONTEND_URL"); vercelURL != "" {
		config.AllowOrigins = append(config.AllowOrigins, vercelURL)
	}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Public routes
	r.POST("/api/register", handlers.Register)
	r.POST("/api/login", handlers.Login)
	r.GET("/:slug", handlers.RedirectURL) // The redirect hot path

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.JWTAuthMiddleware())
	api.Use(middleware.RateLimitMiddleware()) // Protect all API endpoints with Rate Limiting
	{
		api.POST("/shorten", handlers.ShortenURL)
		api.GET("/links", handlers.GetLinks)
		api.GET("/links/:slug/analytics", handlers.GetLinkAnalytics)
		api.GET("/links/:slug/summary", handlers.GetLinkSummary)
	}

	// Start Background Workers
	workers.StartAISummaryWorker()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

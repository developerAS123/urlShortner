package main

import (
	"log"
	"os"

	"github.com/ankitsingh/urlshortener/internal/handlers"
	"github.com/ankitsingh/urlshortener/internal/middleware"
	"github.com/ankitsingh/urlshortener/internal/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (for local development)
	_ = godotenv.Load()

	// Initialize Database and Redis
	repository.InitDB()
	repository.InitRedis()

	r := gin.Default()

	// CORS Setup
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:5173", // Local frontend
		// We'll add Vercel production URL later via ENV or hardcode when known
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
	{
		api.POST("/shorten", handlers.ShortenURL)
		// More analytics endpoints will be added here in Week 2
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

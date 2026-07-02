package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ankitsingh/urlshortener/internal/models"
	"github.com/ankitsingh/urlshortener/internal/repository"
	"github.com/ankitsingh/urlshortener/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ShortenRequest struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
}

func ShortenURL(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL provided"})
		return
	}

	// Generate a unique slug
	var slug string
	for i := 0; i < 5; i++ { // Retry up to 5 times for collision
		slug = utils.GenerateSlug(6)
		var count int64
		repository.DB.Model(&models.Link{}).Where("slug = ?", slug).Count(&count)
		if count == 0 {
			break
		}
		if i == 4 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate unique slug"})
			return
		}
	}

	link := models.Link{
		Slug:        slug,
		OriginalURL: req.OriginalURL,
		UserID:      userID.(uint),
		IsActive:    true,
	}

	if err := repository.DB.Create(&link).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save link"})
		return
	}

	// Cache in Redis for 24 hours
	err := repository.RedisClient.Set(repository.Ctx, "slug:"+slug, req.OriginalURL, 24*time.Hour).Err()
	if err != nil {
		log.Printf("Failed to cache slug in Redis: %v", err)
		// We still return success as DB write succeeded
	}

	c.JSON(http.StatusCreated, gin.H{
		"slug":         link.Slug,
		"original_url": link.OriginalURL,
		"short_url":    c.Request.Host + "/" + link.Slug, // Simplistic representation
	})
}

func RedirectURL(c *gin.Context) {
	slug := c.Param("slug")

	// 1. Check Redis Cache First
	originalURL, err := repository.RedisClient.Get(repository.Ctx, "slug:"+slug).Result()
	if err == redis.Nil {
		// Cache miss, check Postgres
		var link models.Link
		if err := repository.DB.Where("slug = ? AND is_active = ?", slug, true).First(&link).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		originalURL = link.OriginalURL
		// Warm up cache
		repository.RedisClient.Set(repository.Ctx, "slug:"+slug, originalURL, 24*time.Hour)
	} else if err != nil {
		log.Printf("Redis error: %v", err)
		// Fallback to DB if redis fails, but for brevity, we could do it here too
	}

	// 2. Async Click Logging
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()
	referrer := c.Request.Referer()

	go func(s string, i string, ua string, ref string) {
		// Retrieve Link ID. If we got originalURL from cache, we might need LinkID.
		// For week 1, doing a quick select to get LinkID
		var link models.Link
		if err := repository.DB.Select("id").Where("slug = ?", s).First(&link).Error; err != nil {
			log.Printf("Failed to find link for click event logging: %v", err)
			return
		}

		click := models.ClickEvent{
			LinkID:    link.ID,
			IPAddress: i,
			UserAgent: ua,
			Referrer:  ref,
			// Note: Country, City, DeviceType, Browser will be parsed in Week 2
		}
		if err := repository.DB.Create(&click).Error; err != nil {
			log.Printf("Failed to log click event: %v", err)
		}
	}(slug, ip, userAgent, referrer)

	// 3. Redirect
	c.Redirect(http.StatusMovedPermanently, originalURL)
}

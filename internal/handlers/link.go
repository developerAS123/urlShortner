package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ankitsingh/urlshortener/internal/models"
	"github.com/ankitsingh/urlshortener/internal/repository"
	"github.com/ankitsingh/urlshortener/internal/services"
	"github.com/ankitsingh/urlshortener/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/mileusna/useragent"
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
	}

	c.JSON(http.StatusCreated, gin.H{
		"slug":         link.Slug,
		"original_url": link.OriginalURL,
		"short_url":    c.Request.Host + "/" + link.Slug,
	})
}

func RedirectURL(c *gin.Context) {
	slug := c.Param("slug")

	originalURL, err := repository.RedisClient.Get(repository.Ctx, "slug:"+slug).Result()
	if err == redis.Nil {
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
		repository.RedisClient.Set(repository.Ctx, "slug:"+slug, originalURL, 24*time.Hour)
	} else if err != nil {
		log.Printf("Redis error: %v", err)
	}

	// Async Click Logging enriched with GeoIP and UserAgent
	ip := c.ClientIP()
	uaString := c.Request.UserAgent()
	referrer := c.Request.Referer()

	go func(s, i, ua, ref string) {
		var link models.Link
		if err := repository.DB.Select("id").Where("slug = ?", s).First(&link).Error; err != nil {
			log.Printf("Failed to find link for click event logging: %v", err)
			return
		}

		country, city := services.GetLocation(i)
		uaParsed := useragent.Parse(ua)
		
		deviceType := "desktop"
		if uaParsed.Mobile {
			deviceType = "mobile"
		} else if uaParsed.Tablet {
			deviceType = "tablet"
		} else if uaParsed.Bot {
			deviceType = "bot"
		}

		click := models.ClickEvent{
			LinkID:     link.ID,
			IPAddress:  i,
			UserAgent:  ua,
			Referrer:   ref,
			Country:    country,
			City:       city,
			DeviceType: deviceType,
			Browser:    uaParsed.Name,
		}
		if err := repository.DB.Create(&click).Error; err != nil {
			log.Printf("Failed to log click event: %v", err)
		}
	}(slug, ip, uaString, referrer)

	c.Redirect(http.StatusMovedPermanently, originalURL)
}

// GetLinks returns all links for a user along with total clicks
func GetLinks(c *gin.Context) {
	userID, _ := c.Get("userID")

	type LinkResult struct {
		models.Link
		TotalClicks int64 `json:"total_clicks"`
	}

	var results []LinkResult

	// Join with click_events to get count
	err := repository.DB.Model(&models.Link{}).
		Select("links.*, COUNT(click_events.id) as total_clicks").
		Joins("LEFT JOIN click_events ON click_events.link_id = links.id").
		Where("links.user_id = ?", userID).
		Group("links.id").
		Order("links.created_at DESC").
		Find(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch links"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetLinkAnalytics returns aggregated analytics for a specific link
func GetLinkAnalytics(c *gin.Context) {
	userID, _ := c.Get("userID")
	slug := c.Param("slug")

	var link models.Link
	if err := repository.DB.Where("slug = ? AND user_id = ?", slug, userID).First(&link).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found or unauthorized"})
		return
	}

	// For Week 2, returning simple aggregation arrays
	// A more complex aggregation by day/hour is useful for recharts line charts
	// Let's do raw queries or GORM grouping

	// Group by Date (Postgres specific date_trunc, can simplify for SQLite/MySQL but we use Postgres)
	type ClicksByDate struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	var clicksByDate []ClicksByDate
	repository.DB.Raw("SELECT DATE(clicked_at) as date, COUNT(*) as count FROM click_events WHERE link_id = ? GROUP BY DATE(clicked_at) ORDER BY date ASC", link.ID).Scan(&clicksByDate)

	type ClicksByCountry struct {
		Country string `json:"country"`
		Count   int64  `json:"count"`
	}
	var clicksByCountry []ClicksByCountry
	repository.DB.Raw("SELECT country, COUNT(*) as count FROM click_events WHERE link_id = ? GROUP BY country ORDER BY count DESC", link.ID).Scan(&clicksByCountry)

	type ClicksByDevice struct {
		DeviceType string `json:"device_type"`
		Count      int64  `json:"count"`
	}
	var clicksByDevice []ClicksByDevice
	repository.DB.Raw("SELECT device_type, COUNT(*) as count FROM click_events WHERE link_id = ? GROUP BY device_type", link.ID).Scan(&clicksByDevice)

	c.JSON(http.StatusOK, gin.H{
		"clicks_by_date":    clicksByDate,
		"clicks_by_country": clicksByCountry,
		"clicks_by_device":  clicksByDevice,
	})
}

// GetLinkSummary returns the latest AI summary for a link
func GetLinkSummary(c *gin.Context) {
	userID, _ := c.Get("userID")
	slug := c.Param("slug")

	var link models.Link
	if err := repository.DB.Where("slug = ? AND user_id = ?", slug, userID).First(&link).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found or unauthorized"})
		return
	}

	var summary models.AISummary
	if err := repository.DB.Where("link_id = ?", link.ID).Order("week_start DESC").First(&summary).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Summary not available yet for this link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"summary":      summary.SummaryText,
		"generated_at": summary.GeneratedAt,
	})
}

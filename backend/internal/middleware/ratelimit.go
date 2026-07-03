package middleware

import (
	"net/http"

	"github.com/ankitsingh/urlshortener/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
)

var limiter *redis_rate.Limiter

func InitRateLimiter() {
	limiter = redis_rate.NewLimiter(repository.RedisClient)
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if limiter == nil {
			c.Next()
			return
		}

		// Rate limit based on IP address
		ip := c.ClientIP()
		
		// Allow 10 requests per minute
		res, err := limiter.Allow(c, ip, redis_rate.PerMinute(10))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		if res.Allowed == 0 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

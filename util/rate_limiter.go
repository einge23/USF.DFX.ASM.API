package util

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

//returns a rate limiter as a gin Handler Function
func GetRateLimiter() gin.HandlerFunc {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  60,
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		limiterCtx, err := instance.Get(c, c.ClientIP())
		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		if limiterCtx.Reached {
			c.JSON(429, gin.H{
				"error":       "Too Many Requests",
				"retry_after": limiterCtx.Reset - time.Now().Unix(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

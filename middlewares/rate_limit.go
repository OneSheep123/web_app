package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

func RateLimitMiddleware(fillInterval time.Duration, capacity int64) func(ctx *gin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, capacity)
	return func(ctx *gin.Context) {
		if bucket.TakeAvailable(1) != 1 {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 1000,
				"msg":  "请求过于频繁请稍后",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

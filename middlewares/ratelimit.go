package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

// 中间件的使用
func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, cap) // 初始化令牌桶
	return func(c *gin.Context) {
		// 如果拿不到令牌  一个
		if bucket.TakeAvailable(1) == 0 {
			c.String(http.StatusOK, "rate limit.......")
			c.Abort() // 终止函数执行
			return
		}
		// 函数执行完毕
		c.Next()
	}
}

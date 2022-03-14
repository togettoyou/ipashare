package middleware

import (
	"time"

	"ipashare/internal/server/middleware/cache"
	"ipashare/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	logger := log.New("middleware").Named("log").L()
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		cost := time.Since(start)
		statusCode := c.Writer.Status()
		uri := c.Request.RequestURI
		data := []zap.Field{
			// 请求用户的 IP
			zap.String("ip", c.ClientIP()),
			// 请求的 RequestURI
			zap.String("uri", uri),
			// 请求的方法
			zap.String("method", c.Request.Method),
			// http状态码
			zap.Int("statusCode", statusCode),
			// 请求花费时间
			zap.Duration("cost", cost),
		}
		result := cache.GetCode(c)
		if statusCode > 499 {
			logger.Error(result, data...)
		} else if statusCode > 399 {
			logger.Warn(result, data...)
		} else {
			logger.Info(result, data...)
		}
	}
}

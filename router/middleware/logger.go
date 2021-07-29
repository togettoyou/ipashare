package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		cost := time.Since(start)
		statusCode := c.Writer.Status()
		uri := c.Request.RequestURI
		data := []zap.Field{
			// 日志类型
			zap.String("type", "super-signature-request-log"),
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
		result := "请求响应"
		if statusCode > 499 {
			zap.L().Error(result, data...)
		} else if statusCode > 399 {
			zap.L().Warn(result, data...)
		} else {
			zap.L().Info(result, data...)
		}
	}
}

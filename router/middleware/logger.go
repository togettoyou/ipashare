package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"super-signature/handler"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type respLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w respLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w respLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// request 请求 Body
		buf, _ := ioutil.ReadAll(c.Request.Body)
		// 把读过的字节流重新放到 Body
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		bodyLogWriter := &respLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		// 覆盖默认的 ResponseWriter
		c.Writer = bodyLogWriter
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
		if gin.IsDebugging() && !strings.HasPrefix(uri, "/swagger/") {
			// Debug 模式开启
			data = append(data,
				// 请求 Query 中数据
				zap.String("requestQuery", c.Request.URL.RawQuery),
				// 请求 Body 中数据
				zap.String("requestBody", string(buf)),
				// response 数据
				zap.String("responseData", bodyLogWriter.body.String()),
			)
		}
		var resp handler.Response
		var result string
		if bodyLogWriter.body.String() != "" {
			err := json.Unmarshal(bodyLogWriter.body.Bytes(), &resp)
			if err == nil {
				result = fmt.Sprintf("\tresponse code: %d\tresponse msg: %s", resp.Code, resp.Msg)
			}
		}
		if statusCode > 499 {
			zap.L().Error(result, data...)
		} else if statusCode > 399 {
			zap.L().Warn(result, data...)
		} else {
			if resp.Code >= 20000 {
				zap.L().Warn(result, data...)
			} else if resp.Code >= 10000 {
				zap.L().Error(result, data...)
			} else {
				zap.L().Info(result, data...)
			}
		}
	}
}

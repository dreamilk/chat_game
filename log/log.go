package log

import (
	"bytes"
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var defaultLogger *zap.Logger

func init() {
	logger, err := zap.NewDevelopment(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	defaultLogger = logger
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	defaultLogger.Info(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	defaultLogger.Error(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	defaultLogger.Warn(msg, fields...)
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	defaultLogger.Debug(msg, fields...)
}

func GinZap() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 使用 ResponseWriter 的替代方案来捕获响应
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		// 记录请求日志
		defaultLogger.Info("HTTP request",
			zap.String("client_ip", c.ClientIP()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		// 记录响应内容
		defaultLogger.Info("HTTP response", zap.String("body", blw.body.String()))
	}
}

// 添加一个自定义的 ResponseWriter
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

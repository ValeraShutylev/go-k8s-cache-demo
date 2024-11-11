package middleware

import (
    "github.com/gin-gonic/gin"
    "time"
	"log/slog"
)

func RequestLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        t := time.Now()

        c.Next()

        latency := time.Since(t)

		slog.Info(
			"Handle incoming request",
			slog.String("method", c.Request.Method),
			slog.String("request url", c.Request.RequestURI),
		)
		slog.Debug(
			"Handle incoming request",
			slog.String("log", "audit"),
			slog.String("method", c.Request.Method),
			slog.String("request url", c.Request.RequestURI),
			slog.String("request body", c.Request.Proto),
			slog.String("latency", latency.String()),
		)
    }
}

func ResponseLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

        c.Next()

		slog.Info(
			"Response",
			slog.String("log", "audit"),
			slog.Int("status", c.Writer.Status()),
			slog.String("method", c.Request.Method),
			slog.String("request url", c.Request.RequestURI),
		)
    }
}
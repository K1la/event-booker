package router

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/K1la/event-booker/internal/api/handler"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// customLoggerMiddleware логирует HTTP запросы с использованием zlog
func customLoggerMiddleware() ginext.HandlerFunc {
	return func(c *ginext.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Обрабатываем запрос
		c.Next()

		// Логируем после обработки
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		zlog.Logger.Info().
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Str("ip", clientIP).
			Dur("latency", latency).
			Int("size", bodySize).
			Msg("HTTP Request")
	}
}

func New(handler *handler.Handler) *ginext.Engine {
	e := ginext.New("")
	e.Use(ginext.Recovery(), ginext.Logger()) //customLoggerMiddleware())

	// API routes
	api := e.Group("/api/events")
	{
		api.POST("", handler.CreateEvent)
		api.POST("/:id/book", handler.CreateBooking)
		api.POST("/:id/confirm", handler.ConfirmBookingPayment)
		//api.POST("/:id", handler.CancelBooking)

		api.GET("/:id", handler.GetEventByID)
		api.GET("", handler.GetEvents)

	}

	// // Frontend: serve files from ./web
	// e.GET("/", func(c *ginext.Context) {
	// 	http.ServeFile(c.Writer, c.Request, "./web/index.html")
	// })

	// // Serve static files from web directory
	// e.Static("/web", "./web")

	// OLD
	// Frontend: serve files from ./web without conflicting wildcard
	e.NoRoute(func(c *ginext.Context) {
		if c.Request.URL.Path == "/" {
			http.ServeFile(c.Writer, c.Request, "./web/index.html")
			return
		}
		// Serve only files under /web/ directly from disk
		if strings.HasPrefix(c.Request.URL.Path, "/web/") {
			safe := filepath.Clean("." + c.Request.URL.Path)
			http.ServeFile(c.Writer, c.Request, safe)
			return
		}
		c.Status(http.StatusNotFound)
	})

	return e
}

package router

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/K1la/event-booker/internal/api/handler"

	"github.com/wb-go/wbf/ginext"
)

func New(handler *handler.Handler) *ginext.Engine {
	e := ginext.New("")
	e.Use(ginext.Recovery(), ginext.Logger())

	// API routes
	api := e.Group("/api/")
	{
		api.POST("/events", handler.CreateEvent)
		api.POST("/events/:id/book", handler.CreateBooking)
		api.POST("/events/:id/confirm", handler.ConfirmPayment)

		api.GET("/events", handler.GetAllEvents)
		api.GET("/events/:id", handler.GetEvent)

	}

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

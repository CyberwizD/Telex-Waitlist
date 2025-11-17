package routes

import (
	"net/http"

	"github.com/CyberwizD/Telex-Waitlist/internal/config"
	"github.com/CyberwizD/Telex-Waitlist/internal/handlers"
	"github.com/CyberwizD/Telex-Waitlist/internal/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter wires all HTTP endpoints.
func SetupRouter(cfg *config.Config, waitlistHandler *handlers.WaitlistHandler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(middleware.CORS(cfg.AllowedOrigins))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")
	api.POST("/waitlist", waitlistHandler.Submit)
	api.GET("/waitlist", waitlistHandler.List)

	return router
}

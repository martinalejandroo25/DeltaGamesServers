package api

import (
	"time"

	"github.com/deltagames/deltagamesservers/internal/api/v1"
	"github.com/deltagames/deltagamesservers/internal/api/websocket"
	"github.com/deltagames/deltagamesservers/internal/config"
	"github.com/deltagames/deltagamesservers/internal/database"
	"github.com/deltagames/deltagamesservers/internal/gameserver"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, db *database.DB, manager *gameserver.Manager) *gin.Engine {
	r := gin.Default()

	// CORS Configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Handlers
	tmplHandler := v1.NewTemplateHandler(cfg)
	serverHandler := v1.NewServerHandler(db, manager)
	wsHub := websocket.NewHub(manager)
	go wsHub.Run()

	// API Group
	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":    "online",
				"system":    "DeltaGamesServers Engine",
				"timestamp": time.Now(),
			})
		})

		// Templates
		api.GET("/templates", tmplHandler.ListTemplates)

		// Servers
		servers := api.Group("/servers")
		{
			servers.GET("", serverHandler.ListServers)
			servers.POST("", serverHandler.CreateServer)
			servers.GET("/:id", serverHandler.GetServer)
			servers.POST("/:id/start", serverHandler.StartServer)
			servers.POST("/:id/stop", serverHandler.StopServer)
			servers.DELETE("/:id", serverHandler.DeleteServer)
			servers.GET("/:id/logs", serverHandler.GetLogs)
			servers.GET("/:id/ws", func(c *gin.Context) {
				id := c.Param("id")
				wsHub.HandleLogsWS(c, id)
			})
		}
	}

	// Serve React Single Page Application (SPA)
	r.Static("/assets", "web/dist/assets")
	r.StaticFile("/favicon.svg", "web/dist/favicon.svg")
	r.NoRoute(func(c *gin.Context) {
		c.File("web/dist/index.html")
	})

	return r
}

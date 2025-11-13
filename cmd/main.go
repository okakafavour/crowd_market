package main

import (
	"CROWD_MARKET/config"
	"CROWD_MARKET/routes"
	"CROWD_MARKET/services"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 🧩 Initialize configuration and services
	config.InitConfig()
	services.InitUserService()
	services.InitProductService()

	router := gin.Default()

	// 🌍 CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ✅ Health check route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Crowd Market API is running 🚀"})
	})

	// 🧭 Register all other routes
	routes.RegisterRoutes(router)

	// 🌐 Pick Render-assigned port or fallback to 8080 locally
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Starting server on port %s...", port)
	router.Run(":" + port)
}

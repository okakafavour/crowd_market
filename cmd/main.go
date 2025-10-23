package main

import (
	"CROWD_MARKET/config"
	"CROWD_MARKET/routes"
	"CROWD_MARKET/services"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	services.InitUserService()
	services.InitProductService()

	router := gin.Default()
	routes.RegisterRoutes(router)

	router.Run(":8080")
}

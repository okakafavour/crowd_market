package routes

import (
	"CROWD_MARKET/controllers"
	"CROWD_MARKET/middleware"
	"CROWD_MARKET/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "welcome to crowd market",
		})
	})

	router.POST("/register", controllers.RegisterUser)
	router.POST("/login", controllers.LoginUser)
	router.POST("/products", controllers.AddProduct)
	router.GET("/products", controllers.GetAllProducts)
	router.GET("/auth/google/login", controllers.GoogleLogin)
	router.GET("/auth/google/callback", controllers.GoogleCallBack)
	router.PUT("/:id", controllers.UpdateProduct)
	router.DELETE("/:id", controllers.DeleteProduct)

	router.GET("/verify", func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "verification code is required"})
			return
		}

		err := services.VerifyUserEmail(code)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully!"})
	})

	protected := router.Group("/user")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.GET("/profile", controllers.GetProfile)
	}
}

package routes

import (
	"CROWD_MARKET/controllers"
	"CROWD_MARKET/middleware"
	"CROWD_MARKET/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	// --- Auth routes ---
	router.POST("/register", controllers.RegisterUser)
	router.POST("/login", controllers.LoginUser)

	// --- Google Auth routes ---
	router.GET("/auth/google/login", controllers.GoogleLogin)
	router.GET("/auth/google/callback", controllers.GoogleCallBack)
	router.POST("/auth/google/web", controllers.GoogleLoginWeb)

	// --- Email verification ---
	router.GET("/verify", func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Verification code is required"})
			return
		}

		err := services.VerifyUserEmail(code)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully!"})
	})

	// --- Product routes ---
	productRoutes := router.Group("/products")
	{
		productRoutes.POST("/", controllers.AddProduct)
		productRoutes.GET("/", controllers.GetAllProducts)
		productRoutes.GET("/:id", controllers.GetProductByID)
		productRoutes.PUT("/:id", controllers.UpdateProduct)
		productRoutes.DELETE("/:id", controllers.DeleteProduct)
	}

	// --- Protected routes (JWT required) ---
	protected := router.Group("/user")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.GET("/profile", controllers.GetProfile)
	}
}

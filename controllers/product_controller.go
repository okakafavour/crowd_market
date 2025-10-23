package controllers

import (
	"CROWD_MARKET/model"
	"CROWD_MARKET/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ✅ Add new product
func AddProduct(c *gin.Context) {
	name := c.PostForm("name")
	priceStr := c.PostForm("price")
	area := c.PostForm("area")
	description := c.PostForm("description")
	category := c.PostForm("category")

	if name == "" || priceStr == "" || area == "" || description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price value"})
		return
	}

	file, fileHeader, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}
	defer file.Close()

	imageURL, err := services.UploadToCloudinary(file, fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image: " + err.Error()})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDStr.(string))

	now := time.Now()
	product := model.Product{
		UserID:      userID,
		Name:        name,
		Price:       price,
		Area:        area,
		Description: description,
		ImageURL:    imageURL,
		Category:    category,
		CreatedAt:   &now,
	}

	savedProduct, err := services.AddProduct(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product added successfully",
		"product": savedProduct,
	})
}

// ✅ Get all products (only for the logged-in user)
func GetAllProducts(c *gin.Context) {
	userID := c.GetString("user_id")

	products, err := services.GetAllProducts(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Products fetched successfully",
		"products": products,
	})
}

// ✅ Update product (only if owned by user)
func UpdateProduct(c *gin.Context) {
	productID := c.Param("id")
	userID := c.GetString("user_id")

	updateFields := make(map[string]interface{})

	if c.ContentType() == "application/json" {
		if err := c.ShouldBindJSON(&updateFields); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}
	} else {
		name := c.PostForm("name")
		priceStr := c.PostForm("price")
		area := c.PostForm("area")
		description := c.PostForm("description")
		category := c.PostForm("category")

		if name != "" {
			updateFields["name"] = name
		}
		if priceStr != "" {
			price, err := strconv.ParseFloat(priceStr, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price value"})
				return
			}
			updateFields["price"] = price
		}
		if area != "" {
			updateFields["area"] = area
		}
		if description != "" {
			updateFields["description"] = description
		}
		if category != "" {
			updateFields["category"] = category
		}

		// 🧩 If new image provided → delete old one and upload new
		file, fileHeader, err := c.Request.FormFile("image")
		if err == nil {
			defer file.Close()

			// Get existing product to find old image
			oldProduct, err := services.GetProductByID(productID)
			if err == nil && oldProduct.ImageURL != "" {
				_ = services.DeleteFromCloudinary(oldProduct.ImageURL) // ignore error
			}

			imageURL, err := services.UploadToCloudinary(file, fileHeader)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload new image: " + err.Error()})
				return
			}
			updateFields["image_url"] = imageURL
		}
	}

	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	updatedProduct, err := services.UpdateProductByUser(productID, userID, updateFields)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product": updatedProduct,
	})
}

// ✅ Delete product (only if owned by user)
func DeleteProduct(c *gin.Context) {
	productID := c.Param("id")
	userID := c.GetString("user_id")

	product, err := services.GetProductByID(productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Ownership check
	if product.UserID.Hex() != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this product"})
		return
	}

	// Delete image from Cloudinary
	if product.ImageURL != "" {
		_ = services.DeleteFromCloudinary(product.ImageURL)
	}

	err = services.DeleteProduct(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

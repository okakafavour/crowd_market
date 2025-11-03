package handlers

import (
	"CROWD_MARKET/config"
	"CROWD_MARKET/model"
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ✅ Add new product (price report)
func AddProduct(c *gin.Context) {
	name := c.PostForm("name")
	category := c.PostForm("category")
	area := c.PostForm("area")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")

	// Convert price to float64
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price"})
		return
	}

	// --- handle image ---
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file required"})
		return
	}

	os.MkdirAll("uploads", os.ModePerm)
	filePath := filepath.Join("uploads", filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	now := time.Now()

	product := model.Product{
		ID:          primitive.NewObjectID(),
		Name:        name,
		Category:    category,
		Area:        area,
		Description: description,
		ImageURL:    filePath,
		Price:       price, // assign numeric price
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}

	collection := config.DB.Collection("products")
	_, err = collection.InsertOne(context.Background(), product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product added successfully!"})
}

// ✅ Get all products
func GetProducts(c *gin.Context) {
	collection := config.DB.Collection("products")
	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer cur.Close(context.Background())

	var products []model.Product
	for cur.Next(context.Background()) {
		var product model.Product
		if err := cur.Decode(&product); err == nil {
			products = append(products, product)
		}
	}

	c.JSON(http.StatusOK, products)
}

// ✅ Get single product by ID
func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	collection := config.DB.Collection("products")
	var product model.Product

	err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

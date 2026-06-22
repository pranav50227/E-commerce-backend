package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Product Catalog Service is running"})
	})

	// Product routes
	products := r.Group("/api/v1/products")
	{
		products.GET("/", getProducts)
		products.GET("/:productId", getProductById)
		products.POST("/", createProduct)
		products.PUT("/:productId", updateProduct)
		products.DELETE("/:productId", deleteProduct)
		products.GET("/category/:category", getProductsByCategory)
		products.GET("/search", searchProducts)
	}

	log.Println("Product Catalog Service starting on port 8081...")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to start Product Catalog Service: %v", err)
	}
}

func getProducts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get all products"})
}

func getProductById(c *gin.Context) {
	productId := c.Param("productId")
	c.JSON(http.StatusOK, gin.H{"productId": productId, "message": "Get product by ID"})
}

func createProduct(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully"})
}

func updateProduct(c *gin.Context) {
	productId := c.Param("productId")
	c.JSON(http.StatusOK, gin.H{"productId": productId, "message": "Product updated"})
}

func deleteProduct(c *gin.Context) {
	productId := c.Param("productId")
	c.JSON(http.StatusOK, gin.H{"productId": productId, "message": "Product deleted"})
}

func getProductsByCategory(c *gin.Context) {
	category := c.Param("category")
	c.JSON(http.StatusOK, gin.H{"category": category, "message": "Get products by category"})
}

func searchProducts(c *gin.Context) {
	query := c.Query("q")
	c.JSON(http.StatusOK, gin.H{"query": query, "message": "Search results"})
}

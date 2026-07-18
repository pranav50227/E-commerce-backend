package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	_ "embed"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"shared/auth"
	"shared/constants"
	"shared/utils"
)

//go:embed swagger.json
var swaggerJSON []byte

func main() {
	r := gin.Default()

	// Service mapping configuration
	services := map[string]string{
		"/api/v1/auth":      utils.GetEnv("AUTH_SERVICE_URL", constants.AuthServiceFallbackURL),
		"/api/v1/users":     utils.GetEnv("USER_SERVICE_URL", constants.UserServiceFallbackURL),
		"/api/v1/products":  utils.GetEnv("PRODUCT_SERVICE_URL", constants.ProductServiceFallbackURL),
		"/api/v1/inventory": utils.GetEnv("INVENTORY_SERVICE_URL", constants.InventoryServiceFallbackURL),
		"/api/v1/orders":    utils.GetEnv("ORDER_SERVICE_URL", constants.OrderServiceFallbackURL),
		"/api/v1/cart":      utils.GetEnv("CART_SERVICE_URL", constants.CartServiceFallbackURL),
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "API Gateway is running"})
	})

	// Redirect root to swagger
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// Swagger endpoints
	r.GET("/swagger-doc/doc.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json; charset=utf-8", swaggerJSON)
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger-doc/doc.json")))

	// Dynamic routing proxy middleware/handler
	r.Any("/api/v1/*path", func(c *gin.Context) {
		reqPath := c.Request.URL.Path
		reqMethod := c.Request.Method

		// Determine target service
		var targetURL string
		for prefix, target := range services {
			if strings.HasPrefix(reqPath, prefix) {
				targetURL = target
				break
			}
		}

		if targetURL == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found or path not matched by API Gateway"})
			return
		}

		// Auth Bypass Check
		bypassAuth := false
		if reqPath == "/api/v1/auth/login" || reqPath == "/api/v1/auth/register" {
			bypassAuth = true
		} else if strings.HasPrefix(reqPath, "/api/v1/products") && reqMethod == http.MethodGet {
			bypassAuth = true
		}

		var userId string
		if !bypassAuth {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
				c.Abort()
				return
			}
			token := strings.TrimPrefix(authHeader, "Bearer ")
			var err error
			userId, err = auth.VerifyJWT(token, []byte(constants.DefaultJWTSecret))
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
				c.Abort()
				return
			}
		}

		remote, err := url.Parse(targetURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid service URL configuration"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)

		// Custom Director to ensure host headers and query parameters are preserved correctly
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Host = remote.Host
			if userId != "" {
				req.Header.Set("X-User-Id", userId)
			}
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	})

	port := utils.GetEnv("PORT_GATEWAY", "8000")
	log.Printf("API Gateway starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}


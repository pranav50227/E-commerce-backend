package constants

// Shared JWT authentication constants
const (
	DefaultJWTSecret   = "super-secret-key-12345"
	HeaderAuthorization = "Authorization"
	HeaderXUserID       = "X-User-Id"
)

// Default downstream service fallback URLs
const (
	AuthServiceFallbackURL    = "http://localhost:8085"
	UserServiceFallbackURL    = "http://localhost:8080"
	ProductServiceFallbackURL = "http://localhost:8081"
	InventoryServiceFallbackURL = "http://localhost:8082"
	OrderServiceFallbackURL   = "http://localhost:8083"
	CartServiceFallbackURL    = "http://localhost:8084"
)

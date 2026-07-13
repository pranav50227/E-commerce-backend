package config

import "shared/utils"

// ServiceConfig stores standard service configurations loaded from environment
type ServiceConfig struct {
	Port      string
	SecretKey string
}

// LoadDefaultConfig instantiates config settings using default shared parameters
func LoadDefaultConfig() *ServiceConfig {
	return &ServiceConfig{
		Port:      utils.GetEnv("PORT", "8080"),
		SecretKey: utils.GetEnv("JWT_SECRET", "super-secret-key-12345"),
	}
}

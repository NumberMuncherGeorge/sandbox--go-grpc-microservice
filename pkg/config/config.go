// Package config provides configuration management for the rating service.
package config

import (
	"os"
)

// Config holds the configuration for the service
type Config struct {
	GRPCPort string
	HTTPPort string
	MongoURI string
	MongoDB  string
}

// Load loads the configuration from environment variables
func Load() *Config {
	return &Config{
		GRPCPort: getEnv("GRPC_PORT", "50051"),
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:  getEnv("MONGO_DB", "rating_service"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

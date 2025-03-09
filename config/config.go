package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	GitHub   GitHubConfig
	CORS     CORSConfig
	LogLevel string
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port    string
	GinMode string
}

// GitHubConfig holds GitHub API configuration
type GitHubConfig struct {
	Token    string
	Username string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowOrigins []string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	config := &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "release"),
		},
		GitHub: GitHubConfig{
			Token:    getEnv("GITHUB_TOKEN", ""),
			Username: getEnv("GITHUB_USERNAME", ""),
		},
		CORS: CORSConfig{
			AllowOrigins: strings.Split(getEnv("ALLOW_ORIGINS", "*"), ","),
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	// Validate required configuration
	if config.GitHub.Token == "" {
		logrus.Fatal("GITHUB_TOKEN is required")
	}

	if config.GitHub.Username == "" {
		logrus.Fatal("GITHUB_USERNAME is required")
	}

	return config
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

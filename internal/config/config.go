package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds application configuration
type Config struct {
	Port          int
	GCPBucketName string
	GCPBucketURL  string
	APIKeys       []string
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	cfg := &Config{
		Port:          8080,
		GCPBucketName: "taraxa-snapshot",
		GCPBucketURL:  "https://storage.googleapis.com/storage/v1/b/taraxa-snapshot/o",
	}

	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Port = p
		}
	}

	if bucketName := os.Getenv("GCP_BUCKET_NAME"); bucketName != "" {
		cfg.GCPBucketName = bucketName
	}

	if bucketURL := os.Getenv("GCP_BUCKET_URL"); bucketURL != "" {
		cfg.GCPBucketURL = bucketURL
	}

	if apiKeys := os.Getenv("API_KEYS"); apiKeys != "" {
		cfg.APIKeys = strings.Split(apiKeys, ",")
		// Trim whitespace from each key
		for i, key := range cfg.APIKeys {
			cfg.APIKeys[i] = strings.TrimSpace(key)
		}
	}

	return cfg
}

// IsValidAPIKey checks if the provided API key is valid
func (c *Config) IsValidAPIKey(apiKey string) bool {
	if len(c.APIKeys) == 0 {
		return false
	}

	for _, key := range c.APIKeys {
		if key == apiKey && key != "" {
			return true
		}
	}
	return false
}

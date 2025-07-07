package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	Port          int
	GCPBucketName string
	GCPBucketURL  string
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

	return cfg
}

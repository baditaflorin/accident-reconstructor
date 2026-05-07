// Package config loads runtime configuration from environment variables.
package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

// Config contains all runtime settings for the API server.
type Config struct {
	Addr           string `envconfig:"ADDR"`
	PublicBaseURL  string `envconfig:"PUBLIC_BASE_URL"`
	PagesOrigin    string `envconfig:"PAGES_ORIGIN"`
	StorageDir     string `envconfig:"STORAGE_DIR"`
	MaxUploadBytes int64  `envconfig:"MAX_UPLOAD_BYTES"`
	OllamaBaseURL  string `envconfig:"OLLAMA_BASE_URL"`
	OllamaModel    string `envconfig:"OLLAMA_MODEL"`
	Version        string `envconfig:"VERSION"`
	Commit         string `envconfig:"COMMIT"`
}

// Load reads configuration from environment variables with safe defaults.
func Load() (Config, error) {
	viper.SetDefault("ADDR", ":8080")
	viper.SetDefault("PUBLIC_BASE_URL", "http://localhost:8080")
	viper.SetDefault("PAGES_ORIGIN", "https://baditaflorin.github.io")
	viper.SetDefault("STORAGE_DIR", "./tmp/cases")
	viper.SetDefault("MAX_UPLOAD_BYTES", int64(1024*1024*1024))
	viper.SetDefault("OLLAMA_MODEL", "llama3.1")
	viper.SetDefault("VERSION", "0.1.0")
	viper.SetDefault("COMMIT", "dev")

	for _, key := range []string{
		"ADDR",
		"PUBLIC_BASE_URL",
		"PAGES_ORIGIN",
		"STORAGE_DIR",
		"MAX_UPLOAD_BYTES",
		"OLLAMA_BASE_URL",
		"OLLAMA_MODEL",
		"VERSION",
		"COMMIT",
	} {
		if value, ok := os.LookupEnv(key); ok {
			viper.Set(key, value)
		}
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, fmt.Errorf("read env config: %w", err)
	}

	if cfg.Addr == "" {
		cfg.Addr = viper.GetString("ADDR")
	}
	if cfg.PublicBaseURL == "" {
		cfg.PublicBaseURL = viper.GetString("PUBLIC_BASE_URL")
	}
	if cfg.PagesOrigin == "" {
		cfg.PagesOrigin = viper.GetString("PAGES_ORIGIN")
	}
	if cfg.StorageDir == "" {
		cfg.StorageDir = viper.GetString("STORAGE_DIR")
	}
	if cfg.MaxUploadBytes == 0 {
		cfg.MaxUploadBytes = viper.GetInt64("MAX_UPLOAD_BYTES")
	}
	if cfg.OllamaModel == "" {
		cfg.OllamaModel = viper.GetString("OLLAMA_MODEL")
	}
	if cfg.Version == "" {
		cfg.Version = viper.GetString("VERSION")
	}
	if cfg.Commit == "" {
		cfg.Commit = viper.GetString("COMMIT")
	}

	return cfg, nil
}

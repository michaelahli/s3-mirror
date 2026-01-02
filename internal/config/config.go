package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

// StorageConfig defines configuration for a storage endpoint
type StorageConfig struct {
	Type            string `yaml:"type"` // "s3" or "minio"
	Bucket          string `yaml:"bucket"`
	Region          string `yaml:"region"`
	Endpoint        string `yaml:"endpoint"`          // For MinIO
	AccessKeyID     string `yaml:"access_key_id"`     // For MinIO
	SecretAccessKey string `yaml:"secret_access_key"` // For MinIO
	UseSSL          bool   `yaml:"use_ssl"`           // For MinIO
	Prefix          string `yaml:"prefix"`            // Optional prefix filter
}

// Config holds the application configuration
type Config struct {
	Source  StorageConfig `yaml:"source"`
	Target  StorageConfig `yaml:"target"`
	Workers int           `yaml:"workers"`
	DryRun  bool          `yaml:"dry_run"`
	Verbose bool          `yaml:"verbose"`
}

// LoadFromFile loads configuration from a YAML file
func LoadFromFile(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Set defaults
	if cfg.Workers == 0 {
		cfg.Workers = 10
	}

	// Set default types if not specified
	if cfg.Source.Type == "" {
		cfg.Source.Type = "s3"
	}
	if cfg.Target.Type == "" {
		cfg.Target.Type = "s3"
	}

	// Set default region for S3
	if cfg.Source.Type == "s3" && cfg.Source.Region == "" {
		cfg.Source.Region = "us-east-1"
	}
	if cfg.Target.Type == "s3" && cfg.Target.Region == "" {
		cfg.Target.Region = "us-east-1"
	}

	// Set default SSL for MinIO
	if cfg.Source.Type == "minio" && cfg.Source.Endpoint != "" {
		// UseSSL defaults to false for MinIO (can be overridden in YAML)
	}
	if cfg.Target.Type == "minio" && cfg.Target.Endpoint != "" {
		// UseSSL defaults to false for MinIO
	}

	return &cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if err := c.validateStorageConfig("source", c.Source); err != nil {
		return err
	}
	if err := c.validateStorageConfig("target", c.Target); err != nil {
		return err
	}
	if c.Workers < 1 {
		return errors.New("workers must be at least 1")
	}
	if c.Workers > 100 {
		return errors.New("workers cannot exceed 100")
	}
	return nil
}

func (c *Config) validateStorageConfig(name string, cfg StorageConfig) error {
	if cfg.Type != "s3" && cfg.Type != "minio" {
		return errors.New(name + ": type must be 's3' or 'minio'")
	}
	if cfg.Bucket == "" {
		return errors.New(name + ": bucket cannot be empty")
	}
	if cfg.Type == "minio" {
		if cfg.Endpoint == "" {
			return errors.New(name + ": endpoint is required for MinIO")
		}
		if cfg.AccessKeyID == "" {
			return errors.New(name + ": access_key_id is required for MinIO")
		}
		if cfg.SecretAccessKey == "" {
			return errors.New(name + ": secret_access_key is required for MinIO")
		}
	}
	return nil
}

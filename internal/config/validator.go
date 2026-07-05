package config

import (
	"fmt"
	"strings"
)

// Validate validates the configuration
func (c *Config) Validate() error {
	var errors []string

	// Validate Server
	if c.Server.Port == "" {
		errors = append(errors, "server port is required")
	}

	// Validate Database
	if c.Database.Host == "" {
		errors = append(errors, "database host is required")
	}
	if c.Database.Name == "" {
		errors = append(errors, "database name is required")
	}

	// Validate JWT
	if c.JWT.Secret == "" || c.JWT.Secret == "your-super-secret-key-change-this" {
		errors = append(errors, "JWT secret must be set and should not be default value")
	}
	if c.JWT.Expiration <= 0 {
		errors = append(errors, "JWT expiration must be greater than 0")
	}

	// Validate Casbin
	if c.Casbin.ModelPath == "" {
		errors = append(errors, "Casbin model path is required")
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration validation failed:\n  %s", strings.Join(errors, "\n  "))
	}

	return nil
}
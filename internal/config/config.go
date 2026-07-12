// Package config loads process configuration from environment variables.
package config

import "github.com/kelseyhightower/envconfig"

// Config holds all environment-sourced configuration for the API process.
type Config struct {
	DatabaseURL      string `envconfig:"DATABASE_URL" required:"true"`
	PlaidClientID    string `envconfig:"PLAID_CLIENT_ID" required:"true"`
	PlaidSecret      string `envconfig:"PLAID_SECRET" required:"true"`
	PlaidEnv         string `envconfig:"PLAID_ENV" default:"sandbox"`
	AppEncryptionKey string `envconfig:"APP_ENCRYPTION_KEY" required:"true"`
	HTTPAddr         string `envconfig:"HTTP_ADDR" default:":8080"`
	LogLevel         string `envconfig:"LOG_LEVEL" default:"info"`
}

// Load reads configuration from environment variables into a Config.
func Load() (Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return Config{}, err
	}
	return c, nil
}

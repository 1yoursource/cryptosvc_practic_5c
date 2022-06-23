package crypto

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/d7561985/tel"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/joho/godotenv"
	"projects/practic_5course_cesar/internal/storage"
)

type Config struct {
	Listen  string `env:"LISTEN" envDefault:":2200"`
	Project string `env:"PROJECT_NAME" envDefault:":crypto-service"`
	Tel     tel.Config
	Storage storage.Config
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Listen, validation.Required, validation.Length(2, 50)), //nolint:gomnd
		validation.Field(&c.Project, validation.Required),
	)
}

func ReadConfig() (*Config, error) {
	var cfg Config

	err := godotenv.Load()
	if err != nil {
		return &Config{}, fmt.Errorf("loading .env file: %w", err)
	}

	if err := env.Parse(&cfg); err != nil {
		return &Config{}, fmt.Errorf("env, parse: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return &Config{}, fmt.Errorf("config, validate: %w", err)
	}

	return &cfg, nil
}

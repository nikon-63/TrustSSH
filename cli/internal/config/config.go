package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Region                 string `json:"region"`
	CognitoDomain          string `json:"cognito_domain"`
	ClientID               string `json:"client_id"`
	RedirectURI            string `json:"redirect_uri"`
	APIBaseURL             string `json:"api_base_url"`
	DefaultDurationSeconds int    `json:"default_duration_seconds"`
}

func Load() (Config, error) {
	if err := ensureTrustSSHDir(); err != nil {
		return Config{}, err
	}

	path := ConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, fmt.Errorf("missing config file: %s\nCreate it from your Terraform outputs before running trustssh login", path)
		}
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	cfg.CognitoDomain = strings.TrimRight(cfg.CognitoDomain, "/")
	cfg.APIBaseURL = strings.TrimRight(cfg.APIBaseURL, "/")
	return cfg, nil
}

func TrustSSHDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".trustssh"
	}
	return filepath.Join(home, ".trustssh")
}

func ConfigPath() string {
	return filepath.Join(TrustSSHDir(), "config.json")
}

func TokensPath() string {
	return filepath.Join(TrustSSHDir(), "tokens.json")
}

func ensureTrustSSHDir() error {
	if err := os.MkdirAll(TrustSSHDir(), 0700); err != nil {
		return fmt.Errorf("create trustssh directory: %w", err)
	}
	return os.Chmod(TrustSSHDir(), 0700)
}

func (c Config) validate() error {
	var missing []string
	if c.Region == "" {
		missing = append(missing, "region")
	}
	if c.CognitoDomain == "" {
		missing = append(missing, "cognito_domain")
	}
	if c.ClientID == "" {
		missing = append(missing, "client_id")
	}
	if c.RedirectURI == "" {
		missing = append(missing, "redirect_uri")
	}
	if len(missing) > 0 {
		return fmt.Errorf("config missing required field(s): %s", strings.Join(missing, ", "))
	}
	return nil
}

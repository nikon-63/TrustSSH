package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
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

	cfg, err := Parse(data)
	if err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

func FetchRemote(baseURL string) (Config, string, error) {
	configURL, err := remoteConfigURL(baseURL)
	if err != nil {
		return Config{}, "", err
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(configURL)
	if err != nil {
		return Config{}, "", fmt.Errorf("fetch remote config: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return Config{}, "", fmt.Errorf("read remote config: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return Config{}, "", fmt.Errorf("fetch remote config: status %d", resp.StatusCode)
	}

	cfg, err := Parse(body)
	if err != nil {
		return Config{}, "", fmt.Errorf("parse remote config: %w", err)
	}
	return cfg, configURL, nil
}

func Parse(data []byte) (Config, error) {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	if err := cfg.normalizeAndValidate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func Save(cfg Config) error {
	if err := cfg.normalizeAndValidate(); err != nil {
		return err
	}
	if err := ensureTrustSSHDir(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	if err := os.WriteFile(ConfigPath(), append(data, '\n'), 0600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return os.Chmod(ConfigPath(), 0600)
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

func CertificatePath() string {
	return filepath.Join(TrustSSHDir(), "id_ed25519-cert.pub")
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
	if c.APIBaseURL == "" {
		missing = append(missing, "api_base_url")
	}
	if len(missing) > 0 {
		return fmt.Errorf("config missing required field(s): %s", strings.Join(missing, ", "))
	}
	return nil
}

func (c *Config) normalizeAndValidate() error {
	if err := c.validate(); err != nil {
		return err
	}
	c.CognitoDomain = strings.TrimRight(c.CognitoDomain, "/")
	c.APIBaseURL = strings.TrimRight(c.APIBaseURL, "/")
	if c.DefaultDurationSeconds == 0 {
		c.DefaultDurationSeconds = 1800
	}
	return nil
}

func remoteConfigURL(baseURL string) (string, error) {
	if !strings.Contains(baseURL, "://") {
		baseURL = "https://" + baseURL
	}

	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("parse base URL: %w", err)
	}
	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		return "", fmt.Errorf("base URL must use https or http")
	}
	if parsed.Host == "" {
		return "", fmt.Errorf("base URL must include a host")
	}

	parsed.Path = strings.TrimRight(parsed.Path, "/") + "/config.json"
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String(), nil
}

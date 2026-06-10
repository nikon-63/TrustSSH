package config

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigValidateRequiresRequiredFields(t *testing.T) {
	cfg := Config{}
	if err := cfg.validate(); err == nil {
		t.Fatal("validate returned nil for empty config")
	}
}

func TestConfigValidateAllowsLoginConfig(t *testing.T) {
	cfg := Config{
		Region:        "eu-west-2",
		CognitoDomain: "https://example.auth.eu-west-2.amazoncognito.com",
		ClientID:      "client",
		RedirectURI:   "http://localhost:8765/callback",
		APIBaseURL:    "https://trustssh.example.com",
	}
	if err := cfg.validate(); err != nil {
		t.Fatalf("validate returned error: %v", err)
	}
}

func TestParseNormalizesConfig(t *testing.T) {
	cfg, err := Parse([]byte(`{
		"region": "eu-west-2",
		"cognito_domain": "https://example.auth.eu-west-2.amazoncognito.com/",
		"client_id": "client",
		"redirect_uri": "http://localhost:8765/callback",
		"api_base_url": "https://trustssh.example.com/"
	}`))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if cfg.CognitoDomain != "https://example.auth.eu-west-2.amazoncognito.com" {
		t.Fatalf("CognitoDomain = %q", cfg.CognitoDomain)
	}
	if cfg.APIBaseURL != "https://trustssh.example.com" {
		t.Fatalf("APIBaseURL = %q", cfg.APIBaseURL)
	}
	if cfg.DefaultDurationSeconds != 1800 {
		t.Fatalf("DefaultDurationSeconds = %d", cfg.DefaultDurationSeconds)
	}
}

func TestFetchRemoteAddsConfigPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/config.json" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"region": "eu-west-2",
			"cognito_domain": "https://example.auth.eu-west-2.amazoncognito.com",
			"client_id": "client",
			"redirect_uri": "http://localhost:8765/callback",
			"api_base_url": "https://trustssh.example.com",
			"default_duration_seconds": 1800
		}`))
	}))
	defer server.Close()

	cfg, sourceURL, err := FetchRemote(server.URL)
	if err != nil {
		t.Fatalf("FetchRemote returned error: %v", err)
	}
	if sourceURL != server.URL+"/config.json" {
		t.Fatalf("sourceURL = %q", sourceURL)
	}
	if cfg.ClientID != "client" {
		t.Fatalf("ClientID = %q", cfg.ClientID)
	}
}

func TestRemoteConfigURLDefaultsToHTTPS(t *testing.T) {
	got, err := remoteConfigURL("trustssh.example.com")
	if err != nil {
		t.Fatalf("remoteConfigURL returned error: %v", err)
	}
	if got != "https://trustssh.example.com/config.json" {
		t.Fatalf("remoteConfigURL = %q", got)
	}
}

func TestSetDefaultDurationSecondsCreatesConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := SetDefaultDurationSeconds(2700); err != nil {
		t.Fatalf("SetDefaultDurationSeconds returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(home, ".trustssh", "config.json"))
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	var values map[string]int
	if err := json.Unmarshal(data, &values); err != nil {
		t.Fatalf("parse config: %v", err)
	}
	if values["default_duration_seconds"] != 2700 {
		t.Fatalf("default_duration_seconds = %d", values["default_duration_seconds"])
	}

	dirInfo, err := os.Stat(filepath.Join(home, ".trustssh"))
	if err != nil {
		t.Fatalf("stat trustssh dir: %v", err)
	}
	if got := dirInfo.Mode().Perm(); got != 0700 {
		t.Fatalf("trustssh dir mode = %o", got)
	}
	fileInfo, err := os.Stat(filepath.Join(home, ".trustssh", "config.json"))
	if err != nil {
		t.Fatalf("stat config: %v", err)
	}
	if got := fileInfo.Mode().Perm(); got != 0600 {
		t.Fatalf("config mode = %o", got)
	}
}

func TestSetDefaultDurationSecondsPreservesExistingConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if err := os.MkdirAll(filepath.Join(home, ".trustssh"), 0700); err != nil {
		t.Fatalf("create trustssh dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(home, ".trustssh", "config.json"), []byte(`{
		"region": "eu-west-2",
		"cognito_domain": "https://example.auth.eu-west-2.amazoncognito.com",
		"client_id": "client",
		"redirect_uri": "http://localhost:8765/callback",
		"api_base_url": "https://trustssh.example.com",
		"default_duration_seconds": 1800
	}`), 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	if err := SetDefaultDurationSeconds(1200); err != nil {
		t.Fatalf("SetDefaultDurationSeconds returned error: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg.Region != "eu-west-2" {
		t.Fatalf("Region = %q", cfg.Region)
	}
	if cfg.DefaultDurationSeconds != 1200 {
		t.Fatalf("DefaultDurationSeconds = %d", cfg.DefaultDurationSeconds)
	}
}

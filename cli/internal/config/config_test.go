package config

import (
	"net/http"
	"net/http/httptest"
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

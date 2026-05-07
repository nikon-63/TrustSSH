package config

import "testing"

func TestConfigValidateRequiresAuthFields(t *testing.T) {
	cfg := Config{}
	if err := cfg.validate(); err == nil {
		t.Fatal("validate returned nil for empty config")
	}
}

func TestConfigValidateAllowsMinimalPhaseOneConfig(t *testing.T) {
	cfg := Config{
		Region:        "eu-west-2",
		CognitoDomain: "https://example.auth.eu-west-2.amazoncognito.com",
		ClientID:      "client",
		RedirectURI:   "http://localhost:8765/callback",
	}
	if err := cfg.validate(); err != nil {
		t.Fatalf("validate returned error: %v", err)
	}
}

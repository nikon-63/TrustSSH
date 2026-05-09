package auth

import (
	"net/url"
	"testing"

	"github.com/nikon-63/TrustSSH/cli/internal/config"
)

func TestPasskeyAddURL(t *testing.T) {
	cfg := config.Config{
		CognitoDomain: "https://auth.trustssh.example.com",
		ClientID:      "client123",
		RedirectURI:   "http://localhost:8765/callback",
	}

	got, err := passkeyAddURL(cfg)
	if err != nil {
		t.Fatalf("passkeyAddURL returned error: %v", err)
	}

	parsed, err := url.Parse(got)
	if err != nil {
		t.Fatalf("parse passkey URL: %v", err)
	}

	if parsed.Scheme != "https" {
		t.Fatalf("scheme = %q", parsed.Scheme)
	}
	if parsed.Host != "auth.trustssh.example.com" {
		t.Fatalf("host = %q", parsed.Host)
	}
	if parsed.Path != "/passkeys/add" {
		t.Fatalf("path = %q", parsed.Path)
	}

	query := parsed.Query()
	if query.Get("client_id") != cfg.ClientID {
		t.Fatalf("client_id = %q", query.Get("client_id"))
	}
	if query.Get("redirect_uri") != cfg.RedirectURI {
		t.Fatalf("redirect_uri = %q", query.Get("redirect_uri"))
	}
}

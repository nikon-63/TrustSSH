package cmd

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigureWarnsWhenOverwritingExistingLocalConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	originalDelay := configureOverwriteDelay
	configureOverwriteDelay = 0
	t.Cleanup(func() {
		configureOverwriteDelay = originalDelay
	})

	trustSSHDir := filepath.Join(home, ".trustssh")
	if err := os.MkdirAll(trustSSHDir, 0700); err != nil {
		t.Fatalf("create TrustSSH dir: %v", err)
	}
	configPath := filepath.Join(trustSSHDir, "config.json")
	if err := os.WriteFile(configPath, []byte(`{
		"region": "eu-west-2",
		"cognito_domain": "https://old.example.com",
		"client_id": "old",
		"redirect_uri": "http://localhost:8765/callback",
		"api_base_url": "https://old.example.com",
		"default_duration_seconds": 2700,
		"set_default_key": true
	}`), 0600); err != nil {
		t.Fatalf("write existing config: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"region": "eu-west-2",
			"cognito_domain": "https://new.example.com",
			"client_id": "new",
			"redirect_uri": "http://localhost:8765/callback",
			"api_base_url": "https://new.example.com"
		}`))
	}))
	defer server.Close()

	output := captureStdout(t, func() {
		if err := Configure(server.URL); err != nil {
			t.Fatalf("Configure returned error: %v", err)
		}
	})

	if !strings.Contains(output, "Warning: "+configPath+" already exists and will be overwritten") {
		t.Fatalf("expected overwrite warning, got: %q", output)
	}
	if !strings.Contains(output, "default_duration_seconds and set_default_key") {
		t.Fatalf("expected warning to mention local settings, got: %q", output)
	}
	if !strings.Contains(output, "Press Ctrl+C within 0 seconds to cancel.") {
		t.Fatalf("expected cancel instruction, got: %q", output)
	}

	updated, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read updated config: %v", err)
	}
	if !strings.Contains(string(updated), `"client_id": "new"`) {
		t.Fatalf("config was not overwritten with remote config:\n%s", string(updated))
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}
	os.Stdout = writer
	defer func() {
		os.Stdout = original
	}()

	fn()

	if err := writer.Close(); err != nil {
		t.Fatalf("close stdout writer: %v", err)
	}

	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	return string(output)
}

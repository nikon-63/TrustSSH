package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nikon-63/TrustSSH/cli/internal/config"
)

func TestIssueCertificate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/issue-cert" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer token" {
			t.Fatalf("Authorization = %q", got)
		}

		var req IssueCertificateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.PublicKey == "" {
			t.Fatal("public key is empty")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(IssueCertificateResponse{
			Certificate: "ssh-ed25519-cert-v01@openssh.com AAAA trustssh",
			ValidUntil:  "2026-05-06T18:30:00Z",
			Principals:  []string{"ubuntu"},
			Serial:      123,
		})
	}))
	defer server.Close()

	resp, err := IssueCertificate(config.Config{APIBaseURL: server.URL}, "token", IssueCertificateRequest{
		PublicKey:                "ssh-ed25519 AAAA trustssh",
		RequestedDurationSeconds: 1800,
	})
	if err != nil {
		t.Fatalf("IssueCertificate returned error: %v", err)
	}
	if resp.Certificate == "" {
		t.Fatal("certificate is empty")
	}
}

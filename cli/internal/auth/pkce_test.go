package auth

import "testing"

func TestNewPKCE(t *testing.T) {
	pkce, err := NewPKCE()
	if err != nil {
		t.Fatalf("NewPKCE returned error: %v", err)
	}
	if pkce.Verifier == "" {
		t.Fatal("Verifier is empty")
	}
	if pkce.Challenge == "" {
		t.Fatal("Challenge is empty")
	}
	if pkce.Method != "S256" {
		t.Fatalf("Method = %q, want S256", pkce.Method)
	}
	if pkce.Verifier == pkce.Challenge {
		t.Fatal("Verifier and challenge should differ")
	}
}

package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/nikon-63/TrustSSH/cli/internal/config"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type Claims struct {
	Subject  string `json:"sub"`
	Email    string `json:"email"`
	Username string `json:"cognito:username"`
}

func ExchangeCode(cfg config.Config, code, verifier string) (Tokens, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", cfg.ClientID)
	form.Set("code", code)
	form.Set("redirect_uri", cfg.RedirectURI)
	form.Set("code_verifier", verifier)

	req, err := http.NewRequest(http.MethodPost, cfg.CognitoDomain+"/oauth2/token", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return Tokens{}, fmt.Errorf("create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Tokens{}, fmt.Errorf("exchange auth code for tokens: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return Tokens{}, fmt.Errorf("read token response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return Tokens{}, fmt.Errorf("token exchange failed: status %d: %s", resp.StatusCode, string(body))
	}

	var tokens Tokens
	if err := json.Unmarshal(body, &tokens); err != nil {
		return Tokens{}, fmt.Errorf("parse token response: %w", err)
	}
	if tokens.IDToken == "" || tokens.AccessToken == "" {
		return Tokens{}, fmt.Errorf("token response did not include expected access and ID tokens")
	}

	return tokens, nil
}

func SaveTokens(tokens Tokens) error {
	if err := os.MkdirAll(config.TrustSSHDir(), 0700); err != nil {
		return fmt.Errorf("create trustssh directory: %w", err)
	}

	data, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		return fmt.Errorf("encode tokens: %w", err)
	}

	if err := os.WriteFile(config.TokensPath(), append(data, '\n'), 0600); err != nil {
		return fmt.Errorf("write tokens: %w", err)
	}
	return os.Chmod(config.TokensPath(), 0600)
}

func DecodeIDTokenClaims(idToken string) (Claims, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return Claims{}, fmt.Errorf("ID token is not a JWT")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, fmt.Errorf("decode ID token payload: %w", err)
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return Claims{}, fmt.Errorf("parse ID token claims: %w", err)
	}
	return claims, nil
}

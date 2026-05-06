package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"

	"github.com/nikon-63/TrustSSH/cli/internal/config"
)

type LoginResult struct {
	Tokens Tokens
	Claims Claims
}

func Login(cfg config.Config) (LoginResult, error) {
	pkce, err := NewPKCE()
	if err != nil {
		return LoginResult{}, err
	}

	state, err := randomState()
	if err != nil {
		return LoginResult{}, err
	}

	authURL, err := authorizationURL(cfg, pkce, state)
	if err != nil {
		return LoginResult{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	callbackCh := make(chan struct {
		result CallbackResult
		err    error
	}, 1)

	go func() {
		result, err := WaitForCallback(ctx, cfg.RedirectURI)
		callbackCh <- struct {
			result CallbackResult
			err    error
		}{result: result, err: err}
	}()

	time.Sleep(150 * time.Millisecond)
	if err := OpenBrowser(authURL); err != nil {
		fmt.Println(authURL)
		return LoginResult{}, err
	}

	callback := <-callbackCh
	if callback.err != nil {
		return LoginResult{}, callback.err
	}
	if callback.result.State != state {
		return LoginResult{}, fmt.Errorf("callback state did not match login request")
	}

	tokens, err := ExchangeCode(cfg, callback.result.Code, pkce.Verifier)
	if err != nil {
		return LoginResult{}, err
	}

	claims, err := DecodeIDTokenClaims(tokens.IDToken)
	if err != nil {
		return LoginResult{}, err
	}

	if err := SaveTokens(tokens); err != nil {
		return LoginResult{}, err
	}

	return LoginResult{Tokens: tokens, Claims: claims}, nil
}

func authorizationURL(cfg config.Config, pkce PKCE, state string) (string, error) {
	base, err := url.Parse(cfg.CognitoDomain + "/oauth2/authorize")
	if err != nil {
		return "", fmt.Errorf("build authorization URL: %w", err)
	}

	query := base.Query()
	query.Set("client_id", cfg.ClientID)
	query.Set("response_type", "code")
	query.Set("scope", "openid email profile")
	query.Set("redirect_uri", cfg.RedirectURI)
	query.Set("code_challenge", pkce.Challenge)
	query.Set("code_challenge_method", pkce.Method)
	query.Set("state", state)
	base.RawQuery = query.Encode()

	return base.String(), nil
}

func randomState() (string, error) {
	random := make([]byte, 24)
	if _, err := rand.Read(random); err != nil {
		return "", fmt.Errorf("generate OAuth state: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(random), nil
}

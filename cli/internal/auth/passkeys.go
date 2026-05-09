package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/nikon-63/TrustSSH/cli/internal/config"
)

func OpenPasskeyEnrollment(cfg config.Config) (string, error) {
	passkeyURL, err := passkeyAddURL(cfg)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	callbackCh := make(chan struct {
		result string
		err    error
	}, 1)

	go func() {
		result, err := WaitForPasskeyResult(ctx, cfg.RedirectURI)
		callbackCh <- struct {
			result string
			err    error
		}{result: result, err: err}
	}()

	time.Sleep(150 * time.Millisecond)

	if err := OpenBrowser(passkeyURL); err != nil {
		fmt.Println(passkeyURL)
		return "", err
	}

	callback := <-callbackCh
	if callback.err != nil {
		return "", callback.err
	}

	return callback.result, nil
}

func WaitForPasskeyResult(ctx context.Context, redirectURI string) (string, error) {
	parsed, err := url.Parse(redirectURI)
	if err != nil {
		return "", fmt.Errorf("parse redirect URI: %w", err)
	}
	if parsed.Scheme != "http" {
		return "", fmt.Errorf("redirect URI must use http for the localhost callback")
	}
	if parsed.Host == "" || parsed.Path == "" {
		return "", fmt.Errorf("redirect URI must include host and path")
	}

	listener, err := net.Listen("tcp", parsed.Host)
	if err != nil {
		return "", fmt.Errorf("start callback listener on %s: %w", parsed.Host, err)
	}
	defer listener.Close()

	resultCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	mux.HandleFunc(parsed.Path, func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if oauthErr := query.Get("error"); oauthErr != "" {
			errCh <- fmt.Errorf("passkey enrollment failed: %s", oauthErr)
			http.Error(w, "TrustSSH passkey enrollment failed. You can close this window.", http.StatusBadRequest)
			return
		}

		result := query.Get("result")
		if result == "" {
			errCh <- fmt.Errorf("passkey enrollment callback missing result")
			http.Error(w, "TrustSSH passkey enrollment did not return a result. You can close this window.", http.StatusBadRequest)
			return
		}

		resultCh <- result
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if result == "success" {
			fmt.Fprintln(w, "Passkey enrollment complete. You can close this window.")
			return
		}
		fmt.Fprintf(w, "Passkey enrollment result: %s. You can close this window.\n", result)
	})

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case result := <-resultCh:
		shutdown(server)
		return result, nil
	case err := <-errCh:
		shutdown(server)
		return "", err
	case <-ctx.Done():
		shutdown(server)
		return "", ctx.Err()
	}
}

func passkeyAddURL(cfg config.Config) (string, error) {
	base, err := url.Parse(cfg.CognitoDomain + "/passkeys/add")
	if err != nil {
		return "", fmt.Errorf("build passkey URL: %w", err)
	}

	query := base.Query()
	query.Set("client_id", cfg.ClientID)
	query.Set("redirect_uri", cfg.RedirectURI)
	base.RawQuery = query.Encode()

	return base.String(), nil
}

package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

type CallbackResult struct {
	Code  string
	State string
}

func WaitForCallback(ctx context.Context, redirectURI string) (CallbackResult, error) {
	parsed, err := url.Parse(redirectURI)
	if err != nil {
		return CallbackResult{}, fmt.Errorf("parse redirect URI: %w", err)
	}
	if parsed.Scheme != "http" {
		return CallbackResult{}, fmt.Errorf("redirect URI must use http for the localhost callback")
	}
	if parsed.Host == "" || parsed.Path == "" {
		return CallbackResult{}, fmt.Errorf("redirect URI must include host and path")
	}

	listener, err := net.Listen("tcp", parsed.Host)
	if err != nil {
		return CallbackResult{}, fmt.Errorf("start callback listener on %s: %w", parsed.Host, err)
	}
	defer listener.Close()

	resultCh := make(chan CallbackResult, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	mux.HandleFunc(parsed.Path, func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if oauthErr := query.Get("error"); oauthErr != "" {
			errCh <- fmt.Errorf("cognito login failed: %s", oauthErr)
			http.Error(w, "TrustSSH login failed. You can close this window.", http.StatusBadRequest)
			return
		}

		code := query.Get("code")
		state := query.Get("state")
		if code == "" {
			errCh <- fmt.Errorf("callback missing authorization code")
			http.Error(w, "TrustSSH login failed. You can close this window.", http.StatusBadRequest)
			return
		}

		resultCh <- CallbackResult{Code: code, State: state}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "TrustSSH login complete. You can close this window.")
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
		return CallbackResult{}, err
	case <-ctx.Done():
		shutdown(server)
		return CallbackResult{}, ctx.Err()
	}
}

func shutdown(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/nikon-63/TrustSSH/cli/internal/config"
)

type IssueCertificateRequest struct {
	PublicKey                string `json:"public_key"`
	RequestedDurationSeconds int    `json:"requested_duration_seconds"`
}

type IssueCertificateResponse struct {
	Certificate string   `json:"certificate"`
	ValidUntil  string   `json:"valid_until"`
	Principals  []string `json:"principals"`
	Serial      uint64   `json:"serial"`
}

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func IssueCertificate(cfg config.Config, accessToken string, request IssueCertificateRequest) (IssueCertificateResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return IssueCertificateResponse{}, fmt.Errorf("encode certificate request: %w", err)
	}

	endpoint := strings.TrimRight(cfg.APIBaseURL, "/") + "/issue-cert"
	httpReq, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return IssueCertificateResponse{}, fmt.Errorf("create certificate request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+accessToken)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return IssueCertificateResponse{}, fmt.Errorf("request certificate: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return IssueCertificateResponse{}, fmt.Errorf("read certificate response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var apiErr errorResponse
		if err := json.Unmarshal(respBody, &apiErr); err == nil && apiErr.Message != "" {
			return IssueCertificateResponse{}, fmt.Errorf("certificate request failed: %s", apiErr.Message)
		}
		return IssueCertificateResponse{}, fmt.Errorf("certificate request failed: status %d", resp.StatusCode)
	}

	var issueResp IssueCertificateResponse
	if err := json.Unmarshal(respBody, &issueResp); err != nil {
		return IssueCertificateResponse{}, fmt.Errorf("parse certificate response: %w", err)
	}
	if issueResp.Certificate == "" {
		return IssueCertificateResponse{}, fmt.Errorf("certificate response did not include a certificate")
	}

	return issueResp, nil
}

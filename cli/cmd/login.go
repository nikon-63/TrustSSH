package cmd

import (
	"fmt"
	"strings"

	"github.com/nikon-63/TrustSSH/cli/internal/api"
	"github.com/nikon-63/TrustSSH/cli/internal/auth"
	"github.com/nikon-63/TrustSSH/cli/internal/config"
	"github.com/nikon-63/TrustSSH/cli/internal/sshkeys"
)

func Login() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Println("Opening browser for Cognito login...")

	result, err := auth.Login(cfg)
	if err != nil {
		return err
	}

	displayName := result.Claims.Email
	if displayName == "" {
		displayName = result.Claims.Username
	}
	if displayName == "" {
		displayName = result.Claims.Subject
	}

	keyPair, err := sshkeys.EnsureDefaultKeyPair()
	if err != nil {
		return err
	}

	publicKey, err := sshkeys.ReadPublicKey(keyPair)
	if err != nil {
		return err
	}

	fmt.Printf("Requesting %d minute certificate...\n", cfg.DefaultDurationSeconds/60)
	certificate, err := api.IssueCertificate(cfg, result.Tokens.AccessToken, api.IssueCertificateRequest{
		PublicKey:                publicKey,
		RequestedDurationSeconds: cfg.DefaultDurationSeconds,
	})
	if err != nil {
		return err
	}
	if err := sshkeys.SaveCertificate(keyPair, certificate.Certificate); err != nil {
		return err
	}

	fmt.Printf("Authenticated as %s\n", displayName)
	fmt.Printf("Using SSH key: %s\n", keyPair.PublicKeyDisplayPath)
	fmt.Printf("Allowed SSH principals: %s\n", strings.Join(certificate.Principals, ", "))
	fmt.Printf("Certificate saved: %s\n", keyPair.CertificateDisplayPath)
	fmt.Printf("Certificate valid until: %s\n", certificate.ValidUntil)
	fmt.Printf("Tokens saved: %s\n", config.TokensPath())
	fmt.Println("You can now use normal SSH commands.")
	return nil
}

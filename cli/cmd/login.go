package cmd

import (
	"fmt"

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

	fmt.Printf("Authenticated as %s\n", displayName)
	fmt.Printf("Using SSH key: %s\n", keyPair.PublicKeyDisplayPath)
	fmt.Printf("Tokens saved: %s\n", config.TokensPath())
	fmt.Println("Successfully logged in to TrustSSH CLI.")
	return nil
}

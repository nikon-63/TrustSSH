package cmd

import (
	"fmt"
	"strings"

	"github.com/nikon-63/TrustSSH/cli/internal/auth"
	"github.com/nikon-63/TrustSSH/cli/internal/config"
)

func PasskeysAdd() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Println("Opening browser to add a passkey...")
	result, err := auth.OpenPasskeyEnrollment(cfg)
	if err != nil {
		return err
	}

	if strings.EqualFold(result, "success") {
		fmt.Println("Passkey enrollment complete.")
		return nil
	}

	fmt.Printf("Passkey enrollment result: %s\n", result)

	return nil
}

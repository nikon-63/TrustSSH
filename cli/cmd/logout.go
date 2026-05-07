package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/nikon-63/TrustSSH/cli/internal/config"
)

func Logout() error {
	tokensPath := config.TokensPath()
	certPath := config.CertificatePath()

	if err := os.Remove(tokensPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("No local tokens found.")
		} else {
			return fmt.Errorf("remove local tokens: %w", err)
		}
	} else {
		fmt.Printf("Removed local tokens: %s\n", tokensPath)
	}

	if err := os.Remove(certPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("No local SSH certificate found.")
			return nil
		}
		return fmt.Errorf("remove local SSH certificate: %w", err)
	}
	fmt.Printf("Removed local SSH certificate: %s\n", certPath)
	return nil
}

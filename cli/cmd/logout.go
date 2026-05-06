package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/nikon-63/TrustSSH/cli/internal/config"
)

func Logout() error {
	tokensPath := config.TokensPath()

	if err := os.Remove(tokensPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("No local tokens found.")
			return nil
		}
		return fmt.Errorf("remove local tokens: %w", err)
	}

	fmt.Printf("Removed local tokens: %s\n", tokensPath)
	return nil
}

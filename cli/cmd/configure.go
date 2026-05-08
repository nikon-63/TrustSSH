package cmd

import (
	"fmt"

	"github.com/nikon-63/TrustSSH/cli/internal/config"
)

func Configure(baseURL string) error {
	cfg, sourceURL, err := config.FetchRemote(baseURL)
	if err != nil {
		return err
	}

	if err := config.Save(cfg); err != nil {
		return err
	}

	fmt.Printf("Fetched config: %s\n", sourceURL)
	fmt.Printf("Config saved: %s\n", config.ConfigPath())
	return nil
}

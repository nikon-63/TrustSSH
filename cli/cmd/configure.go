package cmd

import (
	"fmt"
	"strconv"

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

func ConfigureDefaultDuration(minutesArg string) error {
	seconds, err := parseMinutesToSeconds(minutesArg)
	if err != nil {
		return fmt.Errorf("invalid default duration minutes %q: %w", minutesArg, err)
	}
	if err := config.SetDefaultDurationSeconds(seconds); err != nil {
		return err
	}

	fmt.Printf("Default certificate duration set to %d minutes (%d seconds)\n", seconds/60, seconds)
	fmt.Printf("Config saved: %s\n", config.ConfigPath())
	return nil
}

func parseMinutesToSeconds(minutesArg string) (int, error) {
	minutes, err := strconv.Atoi(minutesArg)
	if err != nil || minutes <= 0 {
		return 0, fmt.Errorf("must be a positive integer")
	}
	maxInt := int(^uint(0) >> 1)
	if minutes > maxInt/60 {
		return 0, fmt.Errorf("is too large")
	}
	return minutes * 60, nil
}

package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nikon-63/TrustSSH/cli/internal/config"
	"github.com/nikon-63/TrustSSH/cli/internal/sshconfig"
	"github.com/nikon-63/TrustSSH/cli/internal/sshkeys"
)

var configureOverwriteDelay = 5 * time.Second

func Configure(baseURL string) error {
	if localConfigExists() {
		fmt.Printf("Warning: %s already exists and will be overwritten, including any local settings such as default_duration_seconds and set_default_key.\n", config.ConfigPath())
		fmt.Printf("Press Ctrl+C within %d seconds to cancel.\n", int(configureOverwriteDelay/time.Second))
		time.Sleep(configureOverwriteDelay)
	}

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

func localConfigExists() bool {
	_, err := os.Stat(config.ConfigPath())
	return err == nil
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

func ConfigureSetDefaultKey(valueArg string) error {
	enabled, err := parseStrictBool(valueArg)
	if err != nil {
		return fmt.Errorf("invalid set-default-key value %q: %w", valueArg, err)
	}

	if enabled {
		if _, err := sshkeys.EnsureDefaultKeyPair(); err != nil {
			return err
		}
		if err := sshconfig.EnsureDefaultKey(); err != nil {
			return err
		}
	} else if err := sshconfig.RemoveDefaultKey(); err != nil {
		return err
	}

	if err := config.SetDefaultKey(enabled); err != nil {
		return err
	}

	fmt.Printf("Set default key: %t\n", enabled)
	if enabled {
		fmt.Printf("Updated SSH config: %s\n", sshconfig.ConfigPath())
	} else {
		fmt.Printf("Removed TrustSSH default key config from: %s\n", sshconfig.ConfigPath())
	}
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

func parseStrictBool(valueArg string) (bool, error) {
	switch valueArg {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("must be true or false")
	}
}

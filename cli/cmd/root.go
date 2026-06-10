package cmd

import (
	"fmt"
	"strconv"
)

var Version = "dev"

func Execute(args []string) error {
	if len(args) == 0 {
		return usage()
	}

	switch args[0] {
	case "configure":
		if len(args) == 3 && args[1] == "--default-duration" {
			return ConfigureDefaultDuration(args[2])
		}
		if len(args) == 3 && args[1] == "--set-default-key" {
			return ConfigureSetDefaultKey(args[2])
		}
		if len(args) == 2 && args[1] == "--passkey-add" {
			return PasskeysAdd()
		}
		if len(args) == 2 && args[1] != "--default-duration" && args[1] != "--set-default-key" {
			return Configure(args[1])
		}
		return configureUsageError()
	case "login":
		durationSeconds, err := parseLoginDuration(args[1:])
		if err != nil {
			return err
		}
		return Login(durationSeconds)
	case "logout":
		if len(args) != 1 {
			return fmt.Errorf("usage: trustssh logout")
		}
		return Logout()
	case "help", "-h", "--help":
		return usage()
	default:
		return fmt.Errorf("unknown command %q\n\n%s", args[0], usageText())
	}
}

func usage() error {
	fmt.Println(usageText())
	return nil
}

func usageText() string {
	return fmt.Sprintf(`Usage:
  trustssh configure <base-url>
  trustssh configure --default-duration minutes
  trustssh configure --set-default-key true|false
  trustssh configure --passkey-add
  trustssh login [-d|--duration minutes]
  trustssh logout

Version: %s`, Version)
}

func configureUsageError() error {
	return fmt.Errorf("usage: trustssh configure <base-url>\n       trustssh configure --default-duration minutes\n       trustssh configure --set-default-key true|false\n       trustssh configure --passkey-add")
}

func parseLoginDuration(args []string) (int, error) {
	if len(args) == 0 {
		return 0, nil
	}
	if len(args) != 2 {
		return 0, fmt.Errorf("usage: trustssh login [-d|--duration minutes]")
	}
	if args[0] != "-d" && args[0] != "--duration" {
		return 0, fmt.Errorf("usage: trustssh login [-d|--duration minutes]")
	}
	minutes, err := strconv.Atoi(args[1])
	if err != nil || minutes <= 0 {
		return 0, fmt.Errorf("invalid duration minutes %q: must be a positive integer", args[1])
	}
	maxInt := int(^uint(0) >> 1)
	if minutes > maxInt/60 {
		return 0, fmt.Errorf("duration minutes %q is too large", args[1])
	}
	return minutes * 60, nil
}

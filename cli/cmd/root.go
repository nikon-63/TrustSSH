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
		if len(args) != 2 {
			return fmt.Errorf("usage: trustssh configure <base-url>")
		}
		return Configure(args[1])
	case "passkeys":
		if len(args) != 2 || args[1] != "add" {
			return fmt.Errorf("usage: trustssh passkeys add")
		}
		return PasskeysAdd()
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
  trustssh passkeys add
  trustssh login [-d minutes]
  trustssh logout

Version: %s`, Version)
}

func parseLoginDuration(args []string) (int, error) {
	if len(args) == 0 {
		return 0, nil
	}
	if len(args) != 2 {
		return 0, fmt.Errorf("usage: trustssh login [-d minutes]")
	}
	if args[0] != "-d" && args[0] != "--duration" {
		return 0, fmt.Errorf("usage: trustssh login [-d minutes]")
	}
	minutes, err := strconv.Atoi(args[1])
	if err != nil || minutes <= 0 {
		return 0, fmt.Errorf("invalid duration minutes: %s", args[1])
	}
	return minutes * 60, nil
}

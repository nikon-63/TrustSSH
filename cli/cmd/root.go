package cmd

import "fmt"

func Execute(args []string) error {
	if len(args) == 0 {
		return usage()
	}

	switch args[0] {
	case "login":
		if len(args) != 1 {
			return fmt.Errorf("usage: trustssh login")
		}
		return Login()
	case "logout":
		if len(args) != 1 {
			return fmt.Errorf("usage: trustssh logout")
		}
		return Logout()
	case "help", "-h", "--help":
		return usage()
	default:
		return fmt.Errorf("unknown command %q\n\nUsage:\n  trustssh login\n  trustssh logout", args[0])
	}
}

func usage() error {
	fmt.Println("Usage:")
	fmt.Println("  trustssh login")
	fmt.Println("  trustssh logout")
	return nil
}

package cmd

import "fmt"

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
  trustssh login
  trustssh logout

Version: %s`, Version)
}

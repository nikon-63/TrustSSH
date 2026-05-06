package auth

import (
	"fmt"
	"os/exec"
	"runtime"
)

func OpenBrowser(rawURL string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", rawURL)
	case "linux":
		cmd = exec.Command("xdg-open", rawURL)
	default:
		return fmt.Errorf("automatic browser opening is not supported on %s; open this URL manually: %s", runtime.GOOS, rawURL)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("open browser: %w\nOpen this URL manually: %s", err, rawURL)
	}
	return nil
}

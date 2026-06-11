package sshconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	beginMarker = "# BEGIN TrustSSH managed block"
	endMarker   = "# END TrustSSH managed block"
)

func EnsureDefaultKey() error {
	if err := ensureSSHDir(); err != nil {
		return err
	}

	path := ConfigPath()
	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read SSH config: %w", err)
	}

	content := removeManagedBlock(string(existing))
	content = strings.TrimRight(content, "\n")
	if content != "" {
		content += "\n\n"
	}
	content += managedBlock()

	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		return fmt.Errorf("write SSH config: %w", err)
	}
	return os.Chmod(path, 0600)
}

func RemoveDefaultKey() error {
	path := ConfigPath()
	existing, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read SSH config: %w", err)
	}

	content := removeManagedBlock(string(existing))
	content = strings.TrimRight(content, "\n")
	if content != "" {
		content += "\n"
	}

	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		return fmt.Errorf("write SSH config: %w", err)
	}
	return os.Chmod(path, 0600)
}

func ConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".ssh", "config")
	}
	return filepath.Join(home, ".ssh", "config")
}

func ensureSSHDir() error {
	dir := filepath.Dir(ConfigPath())
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create SSH config directory: %w", err)
	}
	return os.Chmod(dir, 0700)
}

func managedBlock() string {
	return strings.Join([]string{
		beginMarker,
		"Host *",
		"    IdentityFile ~/.trustssh/id_ed25519",
		"    IdentityFile ~/.ssh/id_ed25519",
		"    IdentitiesOnly yes",
		"    AddKeysToAgent yes",
		endMarker,
		"",
	}, "\n")
}

func removeManagedBlock(content string) string {
	var cleaned strings.Builder
	remaining := content

	for {
		start := strings.Index(remaining, beginMarker)
		if start == -1 {
			cleaned.WriteString(remaining)
			return cleaned.String()
		}

		cleaned.WriteString(remaining[:start])
		afterStart := remaining[start:]
		end := strings.Index(afterStart, endMarker)
		if end == -1 {
			return cleaned.String()
		}

		next := end + len(endMarker)
		if next < len(afterStart) && afterStart[next] == '\r' {
			next++
		}
		if next < len(afterStart) && afterStart[next] == '\n' {
			next++
		}
		if next < len(afterStart) && afterStart[next] == '\n' {
			next++
		}

		remaining = afterStart[next:]
	}
}

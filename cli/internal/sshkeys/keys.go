package sshkeys

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nikon-63/TrustSSH/cli/internal/config"
)

type KeyPair struct {
	PrivateKeyPath        string
	PublicKeyPath         string
	PrivateKeyDisplayPath string
	PublicKeyDisplayPath  string
}

func EnsureDefaultKeyPair() (KeyPair, error) {
	return EnsureKeyPair(config.TrustSSHDir())
}

func EnsureKeyPair(dir string) (KeyPair, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return KeyPair{}, fmt.Errorf("create trustssh directory: %w", err)
	}
	if err := os.Chmod(dir, 0700); err != nil {
		return KeyPair{}, fmt.Errorf("set trustssh directory permissions: %w", err)
	}

	privateKeyPath := filepath.Join(dir, "id_ed25519")
	publicKeyPath := privateKeyPath + ".pub"

	if _, err := os.Stat(privateKeyPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return KeyPair{}, fmt.Errorf("check SSH private key: %w", err)
		}
		if _, err := os.Stat(publicKeyPath); err == nil {
			return KeyPair{}, fmt.Errorf("SSH private key is missing but public key exists: %s", publicKeyPath)
		} else if !errors.Is(err, os.ErrNotExist) {
			return KeyPair{}, fmt.Errorf("check SSH public key before key generation: %w", err)
		}
		if err := generateEd25519Key(privateKeyPath); err != nil {
			return KeyPair{}, err
		}
	}

	if err := os.Chmod(privateKeyPath, 0600); err != nil {
		return KeyPair{}, fmt.Errorf("set SSH private key permissions: %w", err)
	}

	if _, err := os.Stat(publicKeyPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return KeyPair{}, fmt.Errorf("check SSH public key: %w", err)
		}
		if err := writePublicKeyFromPrivateKey(privateKeyPath, publicKeyPath); err != nil {
			return KeyPair{}, err
		}
	}

	if err := os.Chmod(publicKeyPath, 0644); err != nil {
		return KeyPair{}, fmt.Errorf("set SSH public key permissions: %w", err)
	}

	return KeyPair{
		PrivateKeyPath:        privateKeyPath,
		PublicKeyPath:         publicKeyPath,
		PrivateKeyDisplayPath: displayPath(privateKeyPath),
		PublicKeyDisplayPath:  displayPath(publicKeyPath),
	}, nil
}

func generateEd25519Key(privateKeyPath string) error {
	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-f", privateKeyPath, "-N", "", "-C", "trustssh")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("generate SSH key with ssh-keygen: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

func writePublicKeyFromPrivateKey(privateKeyPath, publicKeyPath string) error {
	cmd := exec.Command("ssh-keygen", "-y", "-f", privateKeyPath)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("derive SSH public key with ssh-keygen: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	publicKey := bytes.TrimSpace(stdout.Bytes())
	if len(publicKey) == 0 {
		return fmt.Errorf("derive SSH public key: ssh-keygen returned empty output")
	}

	return os.WriteFile(publicKeyPath, append(publicKey, '\n'), 0644)
}

func displayPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	if path == home {
		return "~"
	}

	prefix := home + string(filepath.Separator)
	if strings.HasPrefix(path, prefix) {
		return "~" + string(filepath.Separator) + strings.TrimPrefix(path, prefix)
	}

	return path
}

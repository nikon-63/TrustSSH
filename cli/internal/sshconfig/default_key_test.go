package sshconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureDefaultKeyCreatesSSHConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := EnsureDefaultKey(); err != nil {
		t.Fatalf("EnsureDefaultKey returned error: %v", err)
	}

	configPath := filepath.Join(home, ".ssh", "config")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read SSH config: %v", err)
	}
	content := string(data)
	for _, want := range []string{
		beginMarker,
		"Host *",
		"    IdentityFile ~/.trustssh/id_ed25519",
		"    IdentityFile ~/.ssh/id_ed25519",
		"    IdentitiesOnly yes",
		"    AddKeysToAgent yes",
		endMarker,
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("SSH config missing %q:\n%s", want, content)
		}
	}
	assertMode(t, filepath.Join(home, ".ssh"), 0700)
	assertMode(t, configPath, 0600)
}

func TestEnsureDefaultKeyPreservesUserConfigAndReplacesManagedBlock(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	sshDir := filepath.Join(home, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		t.Fatalf("create SSH dir: %v", err)
	}
	configPath := filepath.Join(sshDir, "config")
	existing := `Host github.com
    IdentityFile ~/.ssh/github

# BEGIN TrustSSH managed block
Host *
    IdentityFile ~/.old/key
# END TrustSSH managed block

Host example.com
    User ubuntu
`
	if err := os.WriteFile(configPath, []byte(existing), 0600); err != nil {
		t.Fatalf("write SSH config: %v", err)
	}

	if err := EnsureDefaultKey(); err != nil {
		t.Fatalf("EnsureDefaultKey returned error: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read SSH config: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "Host github.com") {
		t.Fatalf("user config before managed block was not preserved:\n%s", content)
	}
	if !strings.Contains(content, "Host example.com") {
		t.Fatalf("user config after managed block was not preserved:\n%s", content)
	}
	if strings.Contains(content, "~/.old/key") {
		t.Fatalf("old managed block was not replaced:\n%s", content)
	}
	if strings.Count(content, beginMarker) != 1 {
		t.Fatalf("expected one managed block, got:\n%s", content)
	}
}

func TestRemoveDefaultKeyRemovesOnlyManagedBlock(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	sshDir := filepath.Join(home, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		t.Fatalf("create SSH dir: %v", err)
	}
	configPath := filepath.Join(sshDir, "config")
	existing := "Host github.com\n    IdentityFile ~/.ssh/github\n\n" + managedBlock() + "\nHost example.com\n    User ubuntu\n"
	if err := os.WriteFile(configPath, []byte(existing), 0600); err != nil {
		t.Fatalf("write SSH config: %v", err)
	}

	if err := RemoveDefaultKey(); err != nil {
		t.Fatalf("RemoveDefaultKey returned error: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read SSH config: %v", err)
	}
	content := string(data)
	if strings.Contains(content, beginMarker) || strings.Contains(content, endMarker) {
		t.Fatalf("managed block was not removed:\n%s", content)
	}
	if !strings.Contains(content, "Host github.com") || !strings.Contains(content, "Host example.com") {
		t.Fatalf("user config was not preserved:\n%s", content)
	}
	assertMode(t, configPath, 0600)
}

func TestEnsureDefaultKeyRemovesMultipleExistingManagedBlocks(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	sshDir := filepath.Join(home, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		t.Fatalf("create SSH dir: %v", err)
	}
	configPath := filepath.Join(sshDir, "config")
	existing := "Host before\n    User git\n\n" + managedBlock() + "\nHost middle\n    User ubuntu\n\n" + managedBlock() + "\nHost after\n    User ec2-user\n"
	if err := os.WriteFile(configPath, []byte(existing), 0600); err != nil {
		t.Fatalf("write SSH config: %v", err)
	}

	if err := EnsureDefaultKey(); err != nil {
		t.Fatalf("EnsureDefaultKey returned error: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read SSH config: %v", err)
	}
	content := string(data)
	if strings.Count(content, beginMarker) != 1 || strings.Count(content, endMarker) != 1 {
		t.Fatalf("expected one managed block after rewrite:\n%s", content)
	}
	for _, want := range []string{"Host before", "Host middle", "Host after"} {
		if !strings.Contains(content, want) {
			t.Fatalf("user config %q was not preserved:\n%s", want, content)
		}
	}
}

func TestRemoveDefaultKeyDropsUnterminatedManagedBlock(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	sshDir := filepath.Join(home, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		t.Fatalf("create SSH dir: %v", err)
	}
	configPath := filepath.Join(sshDir, "config")
	existing := `Host github.com
    IdentityFile ~/.ssh/github

# BEGIN TrustSSH managed block
Host *
    IdentityFile ~/.trustssh/id_ed25519
`
	if err := os.WriteFile(configPath, []byte(existing), 0600); err != nil {
		t.Fatalf("write SSH config: %v", err)
	}

	if err := RemoveDefaultKey(); err != nil {
		t.Fatalf("RemoveDefaultKey returned error: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read SSH config: %v", err)
	}
	content := string(data)
	if strings.Contains(content, beginMarker) || strings.Contains(content, "IdentityFile ~/.trustssh/id_ed25519") {
		t.Fatalf("unterminated managed block was not removed:\n%s", content)
	}
	if !strings.Contains(content, "Host github.com") {
		t.Fatalf("user config before unterminated block was not preserved:\n%s", content)
	}
}

func TestRemoveDefaultKeyIgnoresMissingConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := RemoveDefaultKey(); err != nil {
		t.Fatalf("RemoveDefaultKey returned error: %v", err)
	}
}

func assertMode(t *testing.T, path string, want os.FileMode) {
	t.Helper()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	if got := info.Mode().Perm(); got != want {
		t.Fatalf("%s mode = %o, want %o", path, got, want)
	}
}

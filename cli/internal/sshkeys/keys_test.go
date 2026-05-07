package sshkeys

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureKeyPairCreatesFilesWithStrictPermissions(t *testing.T) {
	dir := t.TempDir()

	keyPair, err := EnsureKeyPair(dir)
	if err != nil {
		t.Fatalf("EnsureKeyPair returned error: %v", err)
	}

	if keyPair.PrivateKeyPath != filepath.Join(dir, "id_ed25519") {
		t.Fatalf("PrivateKeyPath = %q", keyPair.PrivateKeyPath)
	}
	if keyPair.PublicKeyPath != filepath.Join(dir, "id_ed25519.pub") {
		t.Fatalf("PublicKeyPath = %q", keyPair.PublicKeyPath)
	}
	if keyPair.CertificatePath != filepath.Join(dir, "id_ed25519-cert.pub") {
		t.Fatalf("CertificatePath = %q", keyPair.CertificatePath)
	}

	assertMode(t, dir, 0700)
	assertMode(t, keyPair.PrivateKeyPath, 0600)
	assertMode(t, keyPair.PublicKeyPath, 0644)

	publicKey, err := os.ReadFile(keyPair.PublicKeyPath)
	if err != nil {
		t.Fatalf("read public key: %v", err)
	}
	if len(publicKey) == 0 {
		t.Fatal("public key is empty")
	}
}

func TestSaveCertificate(t *testing.T) {
	dir := t.TempDir()
	keyPair, err := EnsureKeyPair(dir)
	if err != nil {
		t.Fatalf("EnsureKeyPair returned error: %v", err)
	}

	if err := SaveCertificate(keyPair, "ssh-ed25519-cert-v01@openssh.com AAAA trustssh"); err != nil {
		t.Fatalf("SaveCertificate returned error: %v", err)
	}

	assertMode(t, keyPair.CertificatePath, 0644)
	cert, err := os.ReadFile(keyPair.CertificatePath)
	if err != nil {
		t.Fatalf("read cert: %v", err)
	}
	if string(cert) != "ssh-ed25519-cert-v01@openssh.com AAAA trustssh\n" {
		t.Fatalf("certificate contents = %q", string(cert))
	}
}

func TestEnsureKeyPairReusesExistingKey(t *testing.T) {
	dir := t.TempDir()

	first, err := EnsureKeyPair(dir)
	if err != nil {
		t.Fatalf("first EnsureKeyPair returned error: %v", err)
	}

	privateKey, err := os.ReadFile(first.PrivateKeyPath)
	if err != nil {
		t.Fatalf("read private key: %v", err)
	}

	second, err := EnsureKeyPair(dir)
	if err != nil {
		t.Fatalf("second EnsureKeyPair returned error: %v", err)
	}

	reusedPrivateKey, err := os.ReadFile(second.PrivateKeyPath)
	if err != nil {
		t.Fatalf("read reused private key: %v", err)
	}

	if string(privateKey) != string(reusedPrivateKey) {
		t.Fatal("private key was not reused")
	}
}

func TestEnsureKeyPairRefusesOrphanedPublicKey(t *testing.T) {
	dir := t.TempDir()
	publicKeyPath := filepath.Join(dir, "id_ed25519.pub")

	if err := os.WriteFile(publicKeyPath, []byte("ssh-ed25519 AAAA orphan\n"), 0644); err != nil {
		t.Fatalf("write orphaned public key: %v", err)
	}

	if _, err := EnsureKeyPair(dir); err == nil {
		t.Fatal("EnsureKeyPair returned nil error for orphaned public key")
	}
}

func assertMode(t *testing.T, path string, want os.FileMode) {
	t.Helper()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}

	got := info.Mode().Perm()
	if got != want {
		t.Fatalf("%s mode = %v, want %v", path, got, want)
	}
}

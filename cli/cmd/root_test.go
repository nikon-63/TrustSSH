package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestParseLoginDuration(t *testing.T) {
	cases := []struct {
		name        string
		args        []string
		want        int
		wantErr     bool
		errContains string
	}{
		{
			name: "no args",
			args: nil,
			want: 0,
		},
		{
			name: "short flag",
			args: []string{"-d", "30"},
			want: 1800,
		},
		{
			name: "long flag",
			args: []string{"--duration", "15"},
			want: 900,
		},
		{
			name:        "missing value short",
			args:        []string{"-d"},
			wantErr:     true,
			errContains: "usage: trustssh login [-d|--duration minutes]",
		},
		{
			name:        "missing value long",
			args:        []string{"--duration"},
			wantErr:     true,
			errContains: "usage: trustssh login [-d|--duration minutes]",
		},
		{
			name:        "unknown flag",
			args:        []string{"--dur", "10"},
			wantErr:     true,
			errContains: "usage: trustssh login [-d|--duration minutes]",
		},
		{
			name:        "non integer",
			args:        []string{"-d", "abc"},
			wantErr:     true,
			errContains: "invalid duration",
		},
		{
			name:        "zero minutes",
			args:        []string{"-d", "0"},
			wantErr:     true,
			errContains: "invalid duration",
		},
		{
			name:        "negative minutes",
			args:        []string{"-d", "-5"},
			wantErr:     true,
			errContains: "invalid duration",
		},
		{
			name:        "extra args",
			args:        []string{"-d", "10", "extra"},
			wantErr:     true,
			errContains: "usage: trustssh login [-d|--duration minutes]",
		},
		{
			name:        "overflow minutes",
			args:        []string{"-d", strconv.Itoa(int(^uint(0)>>1)/60 + 1)},
			wantErr:     true,
			errContains: "too large",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseLoginDuration(tc.args)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tc.errContains != "" && !strings.Contains(err.Error(), tc.errContains) {
					t.Fatalf("expected error containing %q, got %q", tc.errContains, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("expected %d seconds, got %d", tc.want, got)
			}
		})
	}
}

func TestUsageTextIncludesLongDurationFlag(t *testing.T) {
	if !strings.Contains(usageText(), "trustssh login [-d|--duration minutes]") {
		t.Fatalf("usage text should include both duration flags, got: %q", usageText())
	}
}

func TestParseMinutesToSeconds(t *testing.T) {
	got, err := parseMinutesToSeconds("25")
	if err != nil {
		t.Fatalf("parseMinutesToSeconds returned error: %v", err)
	}
	if got != 1500 {
		t.Fatalf("expected 1500 seconds, got %d", got)
	}

	for _, value := range []string{"0", "-1", "abc", strconv.Itoa(int(^uint(0)>>1)/60 + 1)} {
		t.Run(value, func(t *testing.T) {
			if _, err := parseMinutesToSeconds(value); err == nil {
				t.Fatalf("expected error for %q", value)
			}
		})
	}
}

func TestParseStrictBool(t *testing.T) {
	cases := []struct {
		value string
		want  bool
	}{
		{value: "true", want: true},
		{value: "false", want: false},
	}
	for _, tc := range cases {
		t.Run(tc.value, func(t *testing.T) {
			got, err := parseStrictBool(tc.value)
			if err != nil {
				t.Fatalf("parseStrictBool returned error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("expected %t, got %t", tc.want, got)
			}
		})
	}

	for _, value := range []string{"TRUE", "False", "1", "yes", ""} {
		t.Run(value, func(t *testing.T) {
			if _, err := parseStrictBool(value); err == nil {
				t.Fatalf("expected error for %q", value)
			}
		})
	}
}

func TestUsageTextIncludesConfigureDefaultDurationFlag(t *testing.T) {
	if !strings.Contains(usageText(), "trustssh configure --default-duration minutes") {
		t.Fatalf("usage text should include configure default duration flag, got: %q", usageText())
	}
}

func TestUsageTextIncludesConfigureSetDefaultKeyFlag(t *testing.T) {
	if !strings.Contains(usageText(), "trustssh configure --set-default-key true|false") {
		t.Fatalf("usage text should include configure set default key flag, got: %q", usageText())
	}
}

func TestUsageTextIncludesConfigurePasskeyAddFlag(t *testing.T) {
	if !strings.Contains(usageText(), "trustssh configure --passkey-add") {
		t.Fatalf("usage text should include configure passkey add flag, got: %q", usageText())
	}
	if strings.Contains(usageText(), "trustssh passkeys add") {
		t.Fatalf("usage text should not include old passkeys command, got: %q", usageText())
	}
}

func TestConfigureDefaultDurationRequiresValue(t *testing.T) {
	err := Execute([]string{"configure", "--default-duration"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "trustssh configure --default-duration minutes") {
		t.Fatalf("expected configure usage error, got: %q", err.Error())
	}
}

func TestOldPasskeysCommandIsRemoved(t *testing.T) {
	err := Execute([]string{"passkeys", "add"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), `unknown command "passkeys"`) {
		t.Fatalf("expected unknown command error, got: %q", err.Error())
	}
}

func TestConfigureSetDefaultKeyRequiresValue(t *testing.T) {
	err := Execute([]string{"configure", "--set-default-key"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "trustssh configure --set-default-key true|false") {
		t.Fatalf("expected configure usage error, got: %q", err.Error())
	}
}

func TestConfigureSetDefaultKeyUpdatesConfigAndSSHConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := Execute([]string{"configure", "--set-default-key", "true"}); err != nil {
		t.Fatalf("Execute set-default-key true returned error: %v", err)
	}

	cfgData, err := os.ReadFile(filepath.Join(home, ".trustssh", "config.json"))
	if err != nil {
		t.Fatalf("read TrustSSH config: %v", err)
	}
	var cfg map[string]bool
	if err := json.Unmarshal(cfgData, &cfg); err != nil {
		t.Fatalf("parse TrustSSH config: %v", err)
	}
	if !cfg["set_default_key"] {
		t.Fatal("set_default_key was not saved as true")
	}

	sshConfigPath := filepath.Join(home, ".ssh", "config")
	sshConfig, err := os.ReadFile(sshConfigPath)
	if err != nil {
		t.Fatalf("read SSH config: %v", err)
	}
	if !strings.Contains(string(sshConfig), "IdentityFile ~/.trustssh/id_ed25519") {
		t.Fatalf("SSH config missing TrustSSH IdentityFile:\n%s", string(sshConfig))
	}

	if err := Execute([]string{"configure", "--set-default-key", "false"}); err != nil {
		t.Fatalf("Execute set-default-key false returned error: %v", err)
	}
	cfgData, err = os.ReadFile(filepath.Join(home, ".trustssh", "config.json"))
	if err != nil {
		t.Fatalf("read TrustSSH config after disable: %v", err)
	}
	if err := json.Unmarshal(cfgData, &cfg); err != nil {
		t.Fatalf("parse TrustSSH config after disable: %v", err)
	}
	if cfg["set_default_key"] {
		t.Fatal("set_default_key was not saved as false")
	}
	sshConfig, err = os.ReadFile(sshConfigPath)
	if err != nil {
		t.Fatalf("read SSH config after disable: %v", err)
	}
	if strings.Contains(string(sshConfig), "IdentityFile ~/.trustssh/id_ed25519") {
		t.Fatalf("SSH config still contains TrustSSH IdentityFile:\n%s", string(sshConfig))
	}
}

package cmd

import (
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
			errContains: "usage",
		},
		{
			name:        "missing value long",
			args:        []string{"--duration"},
			wantErr:     true,
			errContains: "usage",
		},
		{
			name:        "unknown flag",
			args:        []string{"--dur", "10"},
			wantErr:     true,
			errContains: "usage",
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
			errContains: "usage",
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

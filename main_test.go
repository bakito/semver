package main

import (
	"testing"

	"github.com/coreos/go-semver/semver"
)

func TestLatestTag(t *testing.T) {
	tests := []struct {
		name      string
		tags      []string
		prefix    string
		version   string
		withV     bool
		wantFound bool
	}{
		{
			name:      "filters other prefixes and invalid versions and picks highest semver",
			tags:      []string{"release-v1.9.0", "release-v1.10.0", "other-v99.0.0", "release-not-a-version", "releasex-v100.0.0"},
			prefix:    "release-",
			version:   "1.10.0",
			withV:     true,
			wantFound: true,
		},
		{
			name:      "supports tags without v",
			tags:      []string{"service-1.2.3", "service-2.0.0"},
			prefix:    "service-",
			version:   "2.0.0",
			wantFound: true,
		},
		{
			name:      "without prefix ignores prefixed tags",
			tags:      []string{"v1.0.0", "foo-v0.0.1", "v1.1.0", "foo-v9.0.0"},
			prefix:    "",
			version:   "1.1.0",
			withV:     true,
			wantFound: true,
		},
		{
			name:      "no matching valid version",
			tags:      []string{"other-v3.0.0", "service-invalid"},
			prefix:    "service-",
			wantFound: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, withV, found := latestTag(test.tags, test.prefix)
			if found != test.wantFound {
				t.Fatalf("latestTag() found = %v, want %v", found, test.wantFound)
			}
			if !found {
				if got != nil || withV {
					t.Fatalf("latestTag() = (%v, %v), want (nil, false) when not found", got, withV)
				}
				return
			}
			if got.String() != test.version {
				t.Errorf("latestTag() version = %q, want %q", got.String(), test.version)
			}
			if withV != test.withV {
				t.Errorf("latestTag() withV = %v, want %v", withV, test.withV)
			}
		})
	}
}

func TestNormalizePrefix(t *testing.T) {
	tests := []struct {
		name    string
		prefix  string
		want    string
		wantErr bool
	}{
		{name: "empty prefix", prefix: "", want: ""},
		{name: "appends dash", prefix: "api", want: "api-"},
		{name: "rejects trailing dash", prefix: "api-", wantErr: true},
		{name: "rejects leading dash", prefix: "-api", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := normalizePrefix(test.prefix)
			if (err != nil) != test.wantErr {
				t.Fatalf("normalizePrefix() error = %v, wantErr %v", err, test.wantErr)
			}
			if got != test.want {
				t.Errorf("normalizePrefix() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestFormatVersion(t *testing.T) {
	tests := []struct {
		name    string
		numeric bool
		withV   bool
		want    string
	}{
		{name: "with v", withV: true, want: "v1.2.3"},
		{name: "without v", want: "1.2.3"},
		{name: "numeric with v", numeric: true, withV: true, want: "1.2.3"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := formatVersion(semver.New("1.2.3"), test.numeric, test.withV); got != test.want {
				t.Errorf("formatVersion() = %q, want %q", got, test.want)
			}
		})
	}
}

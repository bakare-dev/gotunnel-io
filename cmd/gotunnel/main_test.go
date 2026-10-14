package main

import (
	"strings"
	"testing"
)

func TestAssetName(t *testing.T) {
	cases := []struct {
		goos, goarch, want string
		wantErr            bool
	}{
		{"darwin", "arm64", "gotunnel-darwin-arm64", false},
		{"linux", "amd64", "gotunnel-linux-amd64", false},
		{"windows", "amd64", "gotunnel-windows-amd64.exe", false},
		{"windows", "arm64", "", true},
		{"freebsd", "amd64", "", true},
	}
	for _, tc := range cases {
		got, err := assetName(tc.goos, tc.goarch)
		if (err != nil) != tc.wantErr || got != tc.want {
			t.Errorf("assetName(%s, %s) = %q, %v; want %q, error %v", tc.goos, tc.goarch, got, err, tc.want, tc.wantErr)
		}
	}
}

func TestFindChecksum(t *testing.T) {
	sum := strings.Repeat("ab", 32)
	file := sum + "  gotunnel-darwin-arm64\n" + strings.Repeat("cd", 32) + "  gotunnel-linux-amd64\n"

	got, err := findChecksum(strings.NewReader(file), "gotunnel-linux-amd64")
	if err != nil || len(got) != 32 || got[0] != 0xcd {
		t.Fatalf("findChecksum = %x, %v", got, err)
	}

	if _, err := findChecksum(strings.NewReader(file), "gotunnel-windows-amd64.exe"); err == nil {
		t.Fatal("expected an error for a missing entry")
	}

	if _, err := findChecksum(strings.NewReader("zz  gotunnel-linux-amd64\n"), "gotunnel-linux-amd64"); err == nil {
		t.Fatal("expected an error for a malformed entry")
	}
}

func TestLocalBuildVersions(t *testing.T) {
	for v, local := range map[string]bool{
		"v1.1.0":                               false,
		"v1.2.0-rc.1":                          false,
		"":                                     true,
		"(devel)":                              true,
		"v1.1.0+dirty":                         true,
		"v1.1.1-0.20261014120000-abcdef123456": true,
		"v1.1.1-0.20261014120000-abcdef123456+dirty": true,
		"v0.0.0-20261014120000-abcdef123456":         true,
	} {
		if got := isLocalBuild(v); got != local {
			t.Errorf("isLocalBuild(%q) = %v, want %v", v, got, local)
		}
	}
}

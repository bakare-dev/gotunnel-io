package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/minio/selfupdate"
)

var releasesURL = "https://github.com/bakare-dev/gotunnel-io/releases"

var pseudoVersion = regexp.MustCompile(`[-.]\d{14}-[0-9a-f]{12}(\+dirty)?$`)

var httpClient = &http.Client{Timeout: 10 * time.Minute}

func main() {
	exe, err := os.Executable()
	if err != nil {
		fail("could not find where gotunnel is installed: %v", err)
	}

	asset, err := assetName(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		fail("%v", err)
	}

	version := releaseVersion()
	label := version
	base := releasesURL + "/download/" + version
	if version == "" {
		label = "latest"
		base = releasesURL + "/latest/download"
	}

	fmt.Fprintf(os.Stderr, "Setting up GoTunnel %s for %s/%s (first run only)...\n", label, runtime.GOOS, runtime.GOARCH)

	sum, err := expectedChecksum(base+"/checksums.txt", asset)
	if err != nil {
		fail("%v", err)
	}

	body, err := open(base + "/" + asset)
	if err != nil {
		fail("%v", err)
	}
	defer body.Close()

	if err := selfupdate.Apply(body, selfupdate.Options{Checksum: sum}); err != nil {
		if rerr := selfupdate.RollbackError(err); rerr != nil {
			fail("could not install GoTunnel, and restoring the launcher also failed: %v", rerr)
		}
		fail("could not install GoTunnel: %v", err)
	}

	fmt.Fprintf(os.Stderr, "✓ Installed GoTunnel %s at %s\n\n", label, exe)
	run(exe, os.Args[1:])
}

func releaseVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	if isLocalBuild(info.Main.Version) {
		return ""
	}
	return info.Main.Version
}

func isLocalBuild(v string) bool {
	return v == "" || v == "(devel)" || strings.HasSuffix(v, "+dirty") || pseudoVersion.MatchString(v)
}

func assetName(goos, goarch string) (string, error) {
	supported := map[string]bool{
		"darwin/amd64":  true,
		"darwin/arm64":  true,
		"linux/amd64":   true,
		"linux/arm64":   true,
		"windows/amd64": true,
	}
	if !supported[goos+"/"+goarch] {
		return "", fmt.Errorf("there is no GoTunnel build for %s/%s", goos, goarch)
	}
	name := "gotunnel-" + goos + "-" + goarch
	if goos == "windows" {
		name += ".exe"
	}
	return name, nil
}

func expectedChecksum(url, asset string) ([]byte, error) {
	body, err := open(url)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	return findChecksum(body, asset)
}

func findChecksum(r io.Reader, asset string) ([]byte, error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 2 || fields[1] != asset {
			continue
		}
		sum, err := hex.DecodeString(fields[0])
		if err != nil || len(sum) != 32 {
			return nil, fmt.Errorf("checksums.txt has a malformed entry for %s", asset)
		}
		return sum, nil
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("could not read checksums.txt: %w", err)
	}
	return nil, errors.New("checksums.txt has no entry for " + asset)
}

func open(url string) (io.ReadCloser, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not download %s: %w", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("could not download %s: HTTP %d", url, resp.StatusCode)
	}
	return resp.Body, nil
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "gotunnel: "+format+"\n", args...)
	fmt.Fprintln(os.Stderr, "Other ways to install: https://docs.gotunnel.bakaredev.site/INSTALL/")
	os.Exit(1)
}

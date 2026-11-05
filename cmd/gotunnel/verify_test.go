package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"aead.dev/minisign"
)

func TestVerifiedChecksum(t *testing.T) {
	pub, priv, err := minisign.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	sum := strings.Repeat("ab", 32)
	sums := []byte(sum + "  gotunnel-linux-amd64\n")
	files := map[string][]byte{
		"/good":     sums,
		"/good.sig": minisign.Sign(priv, sums),
		"/tampered": []byte(strings.Repeat("cd", 32) + "  gotunnel-linux-amd64\n"),
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, ok := files[r.URL.Path]
		if !ok {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	got, err := verifiedChecksum(pub.String(), srv.URL+"/good", srv.URL+"/good.sig", "gotunnel-linux-amd64")
	if err != nil || hex.EncodeToString(got) != sum {
		t.Fatalf("signed checksums: %x, %v", got, err)
	}
	if _, err := verifiedChecksum(pub.String(), srv.URL+"/tampered", srv.URL+"/good.sig", "gotunnel-linux-amd64"); err == nil {
		t.Fatal("tampered checksums accepted")
	}
	if _, err := verifiedChecksum(pub.String(), srv.URL+"/good", srv.URL+"/missing", "gotunnel-linux-amd64"); err == nil {
		t.Fatal("missing signature accepted")
	}
}

func TestEmbeddedPublicKeyParses(t *testing.T) {
	var key minisign.PublicKey
	if err := key.UnmarshalText([]byte(releasePublicKey)); err != nil {
		t.Fatalf("embedded key: %v", err)
	}
}

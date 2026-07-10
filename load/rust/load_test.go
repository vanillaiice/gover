package load_test

import (
	"os"
	"path/filepath"
	"testing"

	load "github.com/vanillaiice/gover/v3/load/rust"
)

func TestFromFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "Cargo.toml")
	if err := os.WriteFile(file, []byte(`# comment
[workspace]

[package]
name = "core"
# version = "0.0.0"
version = "1.2.3"

[package.metadata]
version = "do-not-read"
`), 0644); err != nil {
		t.Fatal(err)
	}

	version, err := load.FromFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if version != "1.2.3" {
		t.Fatalf("got %q, want %q", version, "1.2.3")
	}
}

func TestFromFileMissingVersion(t *testing.T) {
	file := filepath.Join(t.TempDir(), "Cargo.toml")
	if err := os.WriteFile(file, []byte(`[workspace]
members = []
`), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := load.FromFile(file); err == nil {
		t.Fatal("got nil, want missing package version error")
	}
}

func TestFromFilePanic(t *testing.T) {
	file := filepath.Join(t.TempDir(), "Cargo.toml")
	if err := os.WriteFile(file, []byte(`[package]
name = "core"
version = "1.2.3"
`), 0644); err != nil {
		t.Fatal(err)
	}

	if version := load.FromFilePanic(file); version != "1.2.3" {
		t.Fatalf("got %q, want %q", version, "1.2.3")
	}
}

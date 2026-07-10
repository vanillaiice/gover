package load_test

import (
	"os"
	"path/filepath"
	"testing"

	load "github.com/vanillaiice/gover/v3/load/php"
)

func TestFromFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "composer.json")
	if err := os.WriteFile(file, []byte(`{
  "name": "gover/api",
  "version": "1.2.3"
}`), 0644); err != nil {
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

func TestFromFilePanic(t *testing.T) {
	file := filepath.Join(t.TempDir(), "composer.json")
	if err := os.WriteFile(file, []byte(`{"version":"1.2.3"}`), 0644); err != nil {
		t.Fatal(err)
	}

	if version := load.FromFilePanic(file); version != "1.2.3" {
		t.Fatalf("got %q, want %q", version, "1.2.3")
	}
}

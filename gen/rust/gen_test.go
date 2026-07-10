package gen_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	gen "github.com/vanillaiice/gover/v3/gen/rust"
)

func TestUpdateCargoVersion(t *testing.T) {
	file := filepath.Join(t.TempDir(), "Cargo.toml")
	if err := os.WriteFile(file, []byte(`[workspace]

[package]
name = "core"
version = "1.2.3"

[package.metadata]
version = "do-not-change"
`), 0644); err != nil {
		t.Fatal(err)
	}

	out, err := gen.UpdateCargoVersion(file, "2.0.0")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(out), `version = "2.0.0"`) {
		t.Fatalf("got %q, want package version updated", string(out))
	}
	if !strings.Contains(string(out), `version = "do-not-change"`) {
		t.Fatalf("got %q, want metadata version unchanged", string(out))
	}
}

package gen_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vanillaiice/gover/v3/gen"
	"github.com/vanillaiice/gover/v3/lang"
)

func TestVersionDispatch(t *testing.T) {
	t.Run("go", func(t *testing.T) {
		out, err := gen.Version(lang.Go, &gen.Opts{
			PackageName: "version",
			Version:     "v1.2.3",
		})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(out), `const Version = "v1.2.3"`) {
			t.Fatalf("got %q, want Go version source", string(out))
		}
	})

	t.Run("js", func(t *testing.T) {
		file := writeVersionJSON(t, "package.json", "1.2.3")
		out, err := gen.Version(lang.JS, &gen.Opts{
			OutputFile: file,
			Version:    "2.0.0",
		})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(out), `"version": "2.0.0"`) {
			t.Fatalf("got %q, want JS version update", string(out))
		}
	})

	t.Run("rust", func(t *testing.T) {
		file := filepath.Join(t.TempDir(), "Cargo.toml")
		if err := os.WriteFile(file, []byte(`[package]
name = "core"
version = "1.2.3"
`), 0644); err != nil {
			t.Fatal(err)
		}
		out, err := gen.Version(lang.Rust, &gen.Opts{
			OutputFile: file,
			Version:    "2.0.0",
		})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(out), `version = "2.0.0"`) {
			t.Fatalf("got %q, want Rust version update", string(out))
		}
	})

	t.Run("php", func(t *testing.T) {
		file := writeVersionJSON(t, "composer.json", "1.2.3")
		out, err := gen.Version(lang.PHP, &gen.Opts{
			OutputFile: file,
			Version:    "2.0.0",
		})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(out), `"version": "2.0.0"`) {
			t.Fatalf("got %q, want PHP version update", string(out))
		}
	})

	t.Run("invalid", func(t *testing.T) {
		if _, err := gen.Version("ruby", &gen.Opts{}); err == nil {
			t.Fatal("got nil, want unsupported language error")
		}
	})
}

func writeVersionJSON(t *testing.T, name, version string) string {
	t.Helper()

	file := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(file, []byte(`{
  "name": "gover-test",
  "version": "`+version+`"
}`), 0644); err != nil {
		t.Fatal(err)
	}
	return file
}

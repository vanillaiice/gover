package gen_test

import (
	"testing"

	gen "github.com/vanillaiice/gover/v3/gen/go"
)

func TestVersionFile(t *testing.T) {
	const version = "6.9.420"
	const versionFileName = "version.go"
	const packageName = "gover"

	t.Run("local=true", func(t *testing.T) {
		out, err := gen.VersionFile(packageName, version, true)
		if err != nil {
			t.Fatal(err)
		}

		want := "package gover\n\n// version is the current version of the package.\nconst version = \"6.9.420\"\n"

		if string(out) != want {
			t.Errorf("got %q, want %q", string(out), want)
		}
	})

	t.Run("lcoal=false", func(t *testing.T) {
		out, err := gen.VersionFile(packageName, version, false)
		if err != nil {
			t.Fatal(err)
		}

		want := "package gover\n\n// Version is the current version of the package.\nconst Version = \"6.9.420\"\n"

		if string(out) != want {
			t.Errorf("got %q, want %q", string(out), want)
		}
	})
}

package gen_test

import (
	"os"
	"testing"

	"github.com/vanillaiice/gover/gen"
)

func TestVersionFile(t *testing.T) {
	const version = "6.9.420"
	const versionFileName = "version.go"
	const packageName = "gover"

	err := gen.VersionFile(packageName, version, true, versionFileName)
	if err != nil {
		t.Fatal(err)
	}

	want := "package gover\n\n// version is the current version of the package.\nconst version = \"6.9.420\""

	if content, err := os.ReadFile(versionFileName); err != nil {
		t.Fatal(err)
	} else if string(content) != want {
		t.Errorf("got %q, want %q", string(content), want)
	}

	err = gen.VersionFile(packageName, version, false, versionFileName)
	if err != nil {
		t.Fatal(err)
	}

	want = "package gover\n\n// Version is the current version of the package.\nconst Version = \"6.9.420\""

	if content, err := os.ReadFile(versionFileName); err != nil {
		t.Fatal(err)
	} else if string(content) != want {
		t.Errorf("got %q, want %q", string(content), want)
	}

	if err = os.Remove(versionFileName); err != nil {
		panic(err)
	}
}

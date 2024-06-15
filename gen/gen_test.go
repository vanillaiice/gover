package gen_test

import (
	"os"
	"testing"

	"github.com/vanillaiice/gover/gen"
)

func TestVersionFile(t *testing.T) {
	err := gen.VersionFile("gover", "6.9.420", "version.go")
	if err != nil {
		t.Fatal(err)
	}

	want := "package gover\n\n// version is the current version of the package.\nconst version = \"6.9.420\"\n"

	if content, err := os.ReadFile("version.go"); err != nil {
		t.Fatal(err)
	} else if string(content) != want {
		t.Errorf("got %q, want %q", string(content), want)
	}

	if err = os.Remove("version.go"); err != nil {
		panic(err)
	}
}

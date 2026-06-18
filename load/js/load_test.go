package load_test

import (
	"testing"

	load "github.com/vanillaiice/gover/v3/load/js"
)

func TestFromFile(t *testing.T) {
	version, err := load.FromFile("package_test.json")
	if err != nil {
		t.Fatal(err)
	}

	want := "6.9.420"

	if version != want {
		t.Errorf("got %q, want %q", version, want)
	}
}

func TestFromFilePanic(t *testing.T) {
	version := load.FromFilePanic("package_test.json")

	want := "6.9.420"

	if version != want {
		t.Errorf("got %q, want %q", version, want)
	}
}

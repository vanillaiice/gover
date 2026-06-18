package load_test

import (
	"testing"

	load "github.com/vanillaiice/gover/v3/load/go"
)

var files = [2]string{"version_test_global.go", "version_test_local.go"}

func TestFromFile(t *testing.T) {
	const want = "6.9.420"
	for _, file := range files {
		version, err := load.FromFile(file)
		if err != nil {
			t.Fatal(err)
		}

		if version != want {
			t.Errorf("got %q, want %q", version, want)
		}
	}
}

func TestFromFilePanic(t *testing.T) {
	const want = "6.9.420"
	for _, file := range files {
		version := load.FromFilePanic(file)

		if version != want {
			t.Errorf("got %q, want %q", version, want)
		}
	}
}

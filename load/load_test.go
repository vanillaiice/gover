package load_test

import (
	"testing"

	"github.com/vanillaiice/gover/load"
)

var files = [2]string{"version_test_global.go", "version_test_local.go"}

func TestFromFile(t *testing.T) {
	const want = "6.9.420"
	for _, file := range files {
		versionData, err := load.FromFile(file)
		if err != nil {
			t.Fatal(err)
		}

		if versionData.Version != want {
			t.Errorf("got %q, want %q", versionData.Version, want)
		}
	}
}

func TestFromFilePanic(t *testing.T) {
	const want = "6.9.420"
	for _, file := range files {
		versionData := load.FromFilePanic(file)

		if versionData.Version != want {
			t.Errorf("got %q, want %q", versionData.Version, want)
		}
	}
}

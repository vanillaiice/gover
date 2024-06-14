package load_test

import (
	"testing"

	"github.com/vanillaiice/gover/load"
)

func TestFromFile(t *testing.T) {
	versionData, err := load.FromFile("gover_test.json")
	if err != nil {
		t.Fatal(err)
	}

	want := "6.9.420"

	if versionData.Version != want {
		t.Errorf("got %q, want %q", versionData.Version, want)
	}
}

func TestFromFilePanic(t *testing.T) {
	versionData := load.FromFilePanic("gover_test.json")

	want := "6.9.420"

	if versionData.Version != want {
		t.Errorf("got %q, want %q", versionData.Version, want)
	}
}

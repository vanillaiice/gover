package gen_test

import (
	"encoding/json"
	"os"
	"testing"

	gen "github.com/vanillaiice/gover/v3/gen/js"
)

func TestUpdatePackageVersion(t *testing.T) {
	const (
		filePath = "package_test.json"
		want     = "2.0.0"
	)

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err = os.WriteFile(filePath, fileContent, 0644)
		if err != nil {
			t.Fatal(err)
		}
	})

	out, err := gen.UpdatePackageVersion(filePath, want)
	if err != nil {
		t.Fatal(err)
	}

	var data map[string]any
	err = json.Unmarshal(out, &data)
	if err != nil {
		t.Fatal(err)
	}

	if data["version"] != want {
		t.Errorf("got %q, want %q", data["version"], want)
	}
}

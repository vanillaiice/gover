package gen_test

import (
	"encoding/json"
	"os"
	"path/filepath"
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

func TestUpdatePackageVersionOnlyChangesTopLevelVersion(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "package.json")
	if err := os.WriteFile(filePath, []byte(`{
  "name": "gover-test",
  "config": {
    "version": "do-not-change"
  },
  "version": "1.0.0"
}`), 0644); err != nil {
		t.Fatal(err)
	}

	out, err := gen.UpdatePackageVersion(filePath, "2.0.0")
	if err != nil {
		t.Fatal(err)
	}

	var data struct {
		Version string `json:"version"`
		Config  struct {
			Version string `json:"version"`
		} `json:"config"`
	}
	if err := json.Unmarshal(out, &data); err != nil {
		t.Fatal(err)
	}

	if data.Version != "2.0.0" {
		t.Errorf("got top-level version %q, want %q", data.Version, "2.0.0")
	}

	if data.Config.Version != "do-not-change" {
		t.Errorf("got nested version %q, want %q", data.Config.Version, "do-not-change")
	}
}

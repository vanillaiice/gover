package gen_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	gen "github.com/vanillaiice/gover/v3/gen/php"
)

func TestUpdateComposerVersion(t *testing.T) {
	file := filepath.Join(t.TempDir(), "composer.json")
	if err := os.WriteFile(file, []byte(`{
  "name": "gover/api",
  "extra": {
    "version": "do-not-change"
  },
  "version": "1.2.3"
}`), 0644); err != nil {
		t.Fatal(err)
	}

	out, err := gen.UpdateComposerVersion(file, "2.0.0")
	if err != nil {
		t.Fatal(err)
	}

	var data struct {
		Version string `json:"version"`
		Extra   struct {
			Version string `json:"version"`
		} `json:"extra"`
	}
	if err := json.Unmarshal(out, &data); err != nil {
		t.Fatal(err)
	}
	if data.Version != "2.0.0" {
		t.Fatalf("got %q, want top-level version updated", data.Version)
	}
	if data.Extra.Version != "do-not-change" {
		t.Fatalf("got %q, want nested version unchanged", data.Extra.Version)
	}
}

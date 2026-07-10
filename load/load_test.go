package load_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vanillaiice/gover/v3/lang"
	"github.com/vanillaiice/gover/v3/load"
)

func TestFromFileDispatch(t *testing.T) {
	tests := []struct {
		name    string
		lang    lang.Lang
		content string
		want    string
	}{
		{
			name: "version.go",
			lang: lang.Go,
			content: `package version

const Version = "v1.2.3"
`,
			want: "v1.2.3",
		},
		{
			name:    "package.json",
			lang:    lang.JS,
			content: `{"version":"1.2.3"}`,
			want:    "1.2.3",
		},
		{
			name:    "package.json",
			lang:    lang.TS,
			content: `{"version":"1.2.3"}`,
			want:    "1.2.3",
		},
		{
			name: "Cargo.toml",
			lang: lang.Rust,
			content: `[package]
name = "core"
version = "1.2.3"
`,
			want: "1.2.3",
		},
		{
			name:    "composer.json",
			lang:    lang.PHP,
			content: `{"version":"1.2.3"}`,
			want:    "1.2.3",
		},
	}

	for _, test := range tests {
		t.Run(string(test.lang), func(t *testing.T) {
			file := filepath.Join(t.TempDir(), test.name)
			if err := os.WriteFile(file, []byte(test.content), 0644); err != nil {
				t.Fatal(err)
			}

			got, err := load.FromFile(file, test.lang)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf("got %q, want %q", got, test.want)
			}
		})
	}
}

func TestFromFileInvalidLang(t *testing.T) {
	if _, err := load.FromFile("version.txt", "ruby"); err == nil {
		t.Fatal("got nil, want invalid language error")
	}
}

func TestFromFilePanicDispatch(t *testing.T) {
	file := filepath.Join(t.TempDir(), "composer.json")
	if err := os.WriteFile(file, []byte(`{"version":"1.2.3"}`), 0644); err != nil {
		t.Fatal(err)
	}

	if got := load.FromFilePanic(file, lang.PHP); got != "1.2.3" {
		t.Fatalf("got %q, want %q", got, "1.2.3")
	}
}

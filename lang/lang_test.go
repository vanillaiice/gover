package lang_test

import (
	"testing"

	"github.com/vanillaiice/gover/v3/lang"
)

func TestParseLang(t *testing.T) {
	for _, name := range []string{"go", "js", "ts", "rust", "php", "plain"} {
		if got, err := lang.ParseLang(name); err != nil || string(got) != name {
			t.Fatalf("ParseLang(%q) = %q, %v; want %q, nil", name, got, err, name)
		}
	}

	if _, err := lang.ParseLang("ruby"); err == nil {
		t.Fatal("got nil, want invalid lang error")
	}
}

func TestDefaultVersionFilePath(t *testing.T) {
	tests := map[lang.Lang]string{
		lang.Go:    "version/version.go",
		lang.JS:    "package.json",
		lang.TS:    "package.json",
		lang.Rust:  "Cargo.toml",
		lang.PHP:   "composer.json",
		lang.Plain: "version.txt",
	}

	for l, want := range tests {
		got, err := lang.DefaultVersionFilePath(l)
		if err != nil {
			t.Fatalf("DefaultVersionFilePath(%q): %v", l, err)
		}
		if got != want {
			t.Fatalf("DefaultVersionFilePath(%q) = %q, want %q", l, got, want)
		}
	}

	if _, err := lang.DefaultVersionFilePath("ruby"); err == nil {
		t.Fatal("got nil, want invalid lang error")
	}
}

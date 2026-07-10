package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCmd(t *testing.T) {
	t.Run("test get env vars", func(t *testing.T) {
		currentDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		tempDir := t.TempDir()
		os.Chdir(tempDir)
		t.Cleanup(func() {
			os.Chdir(currentDir)
			os.Unsetenv("GOVER_VERSION_FILE")
		})

		want := "my-custom-version-file.go"

		if err := os.WriteFile(
			filepath.Join(tempDir, ".gover"),
			fmt.Appendf([]byte{}, "GOVER_VERSION_FILE=%s", want),
			0644,
		); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(
			filepath.Join(tempDir, want),
			[]byte(`package main
			const version = "v1.2.3"`),
			0644,
		); err != nil {
			t.Fatal(err)
		}

		if err := Exec([]string{"gover", "get"}); err != nil {
			t.Error(err)
		}

		if got := os.Getenv("GOVER_VERSION_FILE"); got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("test init env vars", func(t *testing.T) {
		currentDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		tempDir := t.TempDir()
		os.Chdir(tempDir)
		t.Cleanup(func() {
			os.Chdir(currentDir)
			os.Unsetenv("GOVER_VERSION_FILE")
			os.Unsetenv("GOVER_PACKAGE_NAME")
		})

		versionFile := "my-custom-version-file.go"
		packageName := "mypkg"

		if err := os.WriteFile(
			filepath.Join(tempDir, ".gover"),
			fmt.Appendf([]byte{}, "GOVER_VERSION_FILE=%s\nGOVER_PACKAGE_NAME=%s", versionFile, packageName),
			0644,
		); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(
			filepath.Join(tempDir, versionFile),
			fmt.Appendf([]byte{}, `package %s
			const Version = v1.2.3`, packageName),
			0644,
		); err != nil {
			t.Fatal(err)
		}

		if err := Exec([]string{"gover", "init", "-F"}); err != nil {
			t.Error(err)
		}

		if got := os.Getenv("GOVER_VERSION_FILE"); got != versionFile {
			t.Errorf("got %q want %q", got, versionFile)
		}

		if got := os.Getenv("GOVER_PACKAGE_NAME"); got != packageName {
			t.Errorf("got %q want %q", got, packageName)
		}
	})

	t.Run("test bump uses default file", func(t *testing.T) {
		currentDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		tempDir := t.TempDir()
		if err := os.Chdir(tempDir); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			os.Chdir(currentDir)
		})

		versionFile := filepath.Join(tempDir, "version", "version.go")
		if err := os.MkdirAll(filepath.Dir(versionFile), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(versionFile, []byte(`package version

const Version = "v1.2.3"
`), 0644); err != nil {
			t.Fatal(err)
		}

		if err := Exec([]string{"gover", "bump", "--patch"}); err != nil {
			t.Fatal(err)
		}

		out, err := os.ReadFile(versionFile)
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(out), `const Version = "v1.2.4"`) {
			t.Errorf("got %q, want bumped version", string(out))
		}
	})

	t.Run("test tag uses default file", func(t *testing.T) {
		currentDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		tempDir := t.TempDir()
		if err := os.Chdir(tempDir); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			os.Chdir(currentDir)
		})

		versionFile := filepath.Join(tempDir, "version", "version.go")
		if err := os.MkdirAll(filepath.Dir(versionFile), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(versionFile, []byte(`package version

const Version = "v1.2.3"
`), 0644); err != nil {
			t.Fatal(err)
		}

		if err := Exec([]string{"gover", "tag", "--command", "printf {{ .Version }}"}); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("test bump recursive bumps supported files", func(t *testing.T) {
		currentDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		tempDir := t.TempDir()
		if err := os.Chdir(tempDir); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			os.Chdir(currentDir)
		})

		goVersionFile := filepath.Join(tempDir, "services", "api", "version", "version.go")
		jsPackageFile := filepath.Join(tempDir, "apps", "web", "package.json")
		vendoredPackageFile := filepath.Join(tempDir, "apps", "web", "node_modules", "dep", "package.json")

		for _, file := range []string{goVersionFile, jsPackageFile, vendoredPackageFile} {
			if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
				t.Fatal(err)
			}
		}

		if err := os.WriteFile(goVersionFile, []byte(`package version

const Version = "v1.2.3"
`), 0644); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(jsPackageFile, []byte(`{
  "name": "web",
  "version": "2.3.4"
}`), 0644); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(vendoredPackageFile, []byte(`{
  "name": "dep",
  "version": "9.9.9"
}`), 0644); err != nil {
			t.Fatal(err)
		}

		if err := Exec([]string{"gover", "bump", "--recursive", "--minor"}); err != nil {
			t.Fatal(err)
		}

		goOut, err := os.ReadFile(goVersionFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(goOut), `const Version = "v1.3.0"`) {
			t.Errorf("got %q, want bumped Go version", string(goOut))
		}

		jsOut, err := os.ReadFile(jsPackageFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(jsOut), `"version": "2.4.0"`) {
			t.Errorf("got %q, want bumped JS version", string(jsOut))
		}

		vendoredOut, err := os.ReadFile(vendoredPackageFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(vendoredOut), `"version": "9.9.9"`) {
			t.Errorf("got %q, want vendored package to be unchanged", string(vendoredOut))
		}
	})

	t.Run("test bump path targets one js package", func(t *testing.T) {
		currentDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		tempDir := t.TempDir()
		if err := os.Chdir(tempDir); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			os.Chdir(currentDir)
		})

		webPackageFile := filepath.Join(tempDir, "apps", "web", "package.json")
		adminPackageFile := filepath.Join(tempDir, "apps", "admin", "package.json")
		for _, file := range []string{webPackageFile, adminPackageFile} {
			if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
				t.Fatal(err)
			}
		}

		if err := os.WriteFile(webPackageFile, []byte(`{
  "name": "web",
  "version": "1.2.3"
}`), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(adminPackageFile, []byte(`{
  "name": "admin",
  "version": "4.5.6"
}`), 0644); err != nil {
			t.Fatal(err)
		}

		if err := Exec([]string{"gover", "bump", "--patch", "apps/web"}); err != nil {
			t.Fatal(err)
		}

		webOut, err := os.ReadFile(webPackageFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(webOut), `"version": "1.2.4"`) {
			t.Errorf("got %q, want selected JS package bumped", string(webOut))
		}

		adminOut, err := os.ReadFile(adminPackageFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(adminOut), `"version": "4.5.6"`) {
			t.Errorf("got %q, want unselected JS package unchanged", string(adminOut))
		}
	})

	t.Run("test bump recursive filters by explicit lang", func(t *testing.T) {
		currentDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		tempDir := t.TempDir()
		if err := os.Chdir(tempDir); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			os.Chdir(currentDir)
		})

		goVersionFile := filepath.Join(tempDir, "services", "api", "version", "version.go")
		jsPackageFile := filepath.Join(tempDir, "apps", "web", "package.json")
		for _, file := range []string{goVersionFile, jsPackageFile} {
			if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
				t.Fatal(err)
			}
		}

		if err := os.WriteFile(goVersionFile, []byte(`package version

const Version = "v1.2.3"
`), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(jsPackageFile, []byte(`{
  "name": "web",
  "version": "2.3.4"
}`), 0644); err != nil {
			t.Fatal(err)
		}

		if err := Exec([]string{"gover", "--lang", "go", "bump", "--recursive", "--minor"}); err != nil {
			t.Fatal(err)
		}

		goOut, err := os.ReadFile(goVersionFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(goOut), `const Version = "v1.3.0"`) {
			t.Errorf("got %q, want Go version bumped", string(goOut))
		}

		jsOut, err := os.ReadFile(jsPackageFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(jsOut), `"version": "2.3.4"`) {
			t.Errorf("got %q, want JS package unchanged", string(jsOut))
		}
	})

	t.Run("test bump recursive scopes to path arguments", func(t *testing.T) {
		currentDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		tempDir := t.TempDir()
		if err := os.Chdir(tempDir); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			os.Chdir(currentDir)
		})

		webPackageFile := filepath.Join(tempDir, "apps", "web", "package.json")
		adminPackageFile := filepath.Join(tempDir, "apps", "admin", "package.json")
		for _, file := range []string{webPackageFile, adminPackageFile} {
			if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
				t.Fatal(err)
			}
		}

		if err := os.WriteFile(webPackageFile, []byte(`{
  "name": "web",
  "version": "1.2.3"
}`), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(adminPackageFile, []byte(`{
  "name": "admin",
  "version": "4.5.6"
}`), 0644); err != nil {
			t.Fatal(err)
		}

		if err := Exec([]string{"gover", "bump", "--recursive", "--patch", "apps/web"}); err != nil {
			t.Fatal(err)
		}

		webOut, err := os.ReadFile(webPackageFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(webOut), `"version": "1.2.4"`) {
			t.Errorf("got %q, want scoped JS package bumped", string(webOut))
		}

		adminOut, err := os.ReadFile(adminPackageFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(adminOut), `"version": "4.5.6"`) {
			t.Errorf("got %q, want package outside recursive scope unchanged", string(adminOut))
		}
	})

	t.Run("test bump set writes exact version", func(t *testing.T) {
		currentDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		tempDir := t.TempDir()
		if err := os.Chdir(tempDir); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			os.Chdir(currentDir)
		})

		versionFile := filepath.Join(tempDir, "version", "version.go")
		if err := os.MkdirAll(filepath.Dir(versionFile), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(versionFile, []byte(`package version

const Version = "v1.2.3"
`), 0644); err != nil {
			t.Fatal(err)
		}

		if err := Exec([]string{"gover", "bump", "--set", "v2.0.0"}); err != nil {
			t.Fatal(err)
		}

		out, err := os.ReadFile(versionFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(out), `const Version = "v2.0.0"`) {
			t.Errorf("got %q, want exact version", string(out))
		}
	})

	t.Run("test bump dry run does not write", func(t *testing.T) {
		currentDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		tempDir := t.TempDir()
		if err := os.Chdir(tempDir); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			os.Chdir(currentDir)
		})

		packageFile := filepath.Join(tempDir, "package.json")
		original := []byte(`{
  "name": "web",
  "version": "1.2.3"
}`)
		if err := os.WriteFile(packageFile, original, 0644); err != nil {
			t.Fatal(err)
		}

		if err := Exec([]string{"gover", "--lang", "js", "bump", "--dry-run", "--patch"}); err != nil {
			t.Fatal(err)
		}

		out, err := os.ReadFile(packageFile)
		if err != nil {
			t.Fatal(err)
		}
		if string(out) != string(original) {
			t.Errorf("got %q, want dry run to preserve file", string(out))
		}
	})

	t.Run("test bump recursive supports rust and php", func(t *testing.T) {
		currentDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		tempDir := t.TempDir()
		if err := os.Chdir(tempDir); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			os.Chdir(currentDir)
		})

		cargoFile := filepath.Join(tempDir, "crates", "core", "Cargo.toml")
		composerFile := filepath.Join(tempDir, "packages", "api", "composer.json")
		for _, file := range []string{cargoFile, composerFile} {
			if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
				t.Fatal(err)
			}
		}

		if err := os.WriteFile(cargoFile, []byte(`[package]
name = "core"
version = "1.2.3"

[package.metadata]
version = "do-not-change"
`), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(composerFile, []byte(`{
  "name": "gover/api",
  "version": "2.3.4"
}`), 0644); err != nil {
			t.Fatal(err)
		}

		if err := Exec([]string{"gover", "bump", "--recursive", "--patch"}); err != nil {
			t.Fatal(err)
		}

		cargoOut, err := os.ReadFile(cargoFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(cargoOut), `version = "1.2.4"`) {
			t.Errorf("got %q, want bumped Cargo version", string(cargoOut))
		}
		if !strings.Contains(string(cargoOut), `version = "do-not-change"`) {
			t.Errorf("got %q, want metadata version unchanged", string(cargoOut))
		}

		composerOut, err := os.ReadFile(composerFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(composerOut), `"version": "2.3.5"`) {
			t.Errorf("got %q, want bumped composer version", string(composerOut))
		}
	})

	t.Run("test get json output", func(t *testing.T) {
		unsetEnv(t, "GOVER_VERSION_FILE")

		currentDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		tempDir := t.TempDir()
		if err := os.Chdir(tempDir); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			os.Chdir(currentDir)
		})

		if err := os.WriteFile("Cargo.toml", []byte(`[package]
name = "core"
version = "1.2.3"
`), 0644); err != nil {
			t.Fatal(err)
		}

		out := captureStdout(t, func() {
			if err := Exec([]string{"gover", "--lang", "rust", "get", "--json", "--file", "Cargo.toml"}); err != nil {
				t.Fatal(err)
			}
		})

		if !strings.Contains(out, `"version": "1.2.3"`) {
			t.Errorf("got %q, want json version", out)
		}
		if !strings.Contains(out, `"lang": "rust"`) {
			t.Errorf("got %q, want json lang", out)
		}
	})

	t.Run("test release dry run does not write", func(t *testing.T) {
		unsetEnv(t, "GOVER_VERSION_FILE")

		currentDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		tempDir := t.TempDir()
		if err := os.Chdir(tempDir); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			os.Chdir(currentDir)
		})

		versionFile := filepath.Join(tempDir, "version", "version.go")
		if err := os.MkdirAll(filepath.Dir(versionFile), 0755); err != nil {
			t.Fatal(err)
		}
		original := []byte(`package version

const Version = "v1.2.3"
`)
		if err := os.WriteFile(versionFile, original, 0644); err != nil {
			t.Fatal(err)
		}

		out := captureStdout(t, func() {
			if err := Exec([]string{"gover", "release", "--dry-run", "--patch", "--file", "version/version.go", "--commit-command", "printf commit-{{ .Version }}", "--tag-command", "printf tag-{{ .Version }}"}); err != nil {
				t.Fatal(err)
			}
		})

		if !strings.Contains(out, `v1.2.3 -> v1.2.4`) {
			t.Errorf("got %q, want dry-run bump output", out)
		}

		fileOut, err := os.ReadFile(versionFile)
		if err != nil {
			t.Fatal(err)
		}
		if string(fileOut) != string(original) {
			t.Errorf("got %q, want release dry run to preserve file", string(fileOut))
		}
	})

	t.Run("test check fails invalid version", func(t *testing.T) {
		unsetEnv(t, "GOVER_VERSION_FILE")

		currentDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		tempDir := t.TempDir()
		if err := os.Chdir(tempDir); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			os.Chdir(currentDir)
		})

		if err := os.WriteFile("composer.json", []byte(`{
  "name": "gover/api",
  "version": "nope"
}`), 0644); err != nil {
			t.Fatal(err)
		}

		err = Exec([]string{"gover", "--lang", "php", "check", "--file", "composer.json"})
		if err == nil {
			t.Fatal("got nil, want invalid version error")
		}
	})
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()

	value, ok := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if ok {
			os.Setenv(key, value)
		} else {
			os.Unsetenv(key)
		}
	})
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	stdout := os.Stdout
	read, write, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = write
	t.Cleanup(func() {
		os.Stdout = stdout
	})

	fn()

	if err := write.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = stdout

	out, err := io.ReadAll(read)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

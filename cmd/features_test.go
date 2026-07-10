package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFeatureCommands(t *testing.T) {
	t.Run("bump json set returns structured result", func(t *testing.T) {
		tempDir := chdirTemp(t)
		packageFile := filepath.Join(tempDir, "composer.json")
		if err := os.WriteFile(packageFile, []byte(`{
  "name": "gover/api",
  "version": "1.2.3"
}`), 0644); err != nil {
			t.Fatal(err)
		}

		out := captureStdout(t, func() {
			if err := Exec([]string{"gover", "--lang", "php", "bump", "--json", "--set", "2.0.0", "--file", packageFile}); err != nil {
				t.Fatal(err)
			}
		})

		var results []bumpResult
		if err := json.Unmarshal([]byte(out), &results); err != nil {
			t.Fatal(err)
		}
		if len(results) != 1 {
			t.Fatalf("got %d results, want 1", len(results))
		}
		if results[0].OldVersion != "1.2.3" || results[0].NewVersion != "2.0.0" || results[0].Lang != "php" {
			t.Fatalf("got %+v, want php 1.2.3 -> 2.0.0", results[0])
		}

		content, err := os.ReadFile(packageFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(content), `"version": "2.0.0"`) {
			t.Fatalf("got %q, want file updated", string(content))
		}
	})

	t.Run("bump rejects multiple version actions", func(t *testing.T) {
		tempDir := chdirTemp(t)
		versionFile := filepath.Join(tempDir, "version", "version.go")
		writeGoVersionFile(t, versionFile, "v1.2.3")

		err := Exec([]string{"gover", "bump", "--major", "--patch", "--file", versionFile})
		if err == nil {
			t.Fatal("got nil, want multiple version action error")
		}
		if !strings.Contains(err.Error(), "specify exactly one") {
			t.Fatalf("got %q, want exact-one error", err)
		}
	})

	t.Run("commit dry run json renders command", func(t *testing.T) {
		tempDir := chdirTemp(t)
		versionFile := filepath.Join(tempDir, "version", "version.go")
		writeGoVersionFile(t, versionFile, "v1.2.3")

		out := captureStdout(t, func() {
			if err := Exec([]string{"gover", "commit", "--dry-run", "--json", "--file", versionFile, "--command", "printf commit-{{ .Version }}-{{ .File }}"}); err != nil {
				t.Fatal(err)
			}
		})

		var result commandResult
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatal(err)
		}
		if !result.DryRun || result.Version != "v1.2.3" || !strings.Contains(result.Command, "commit-v1.2.3-") {
			t.Fatalf("got %+v, want dry-run commit command", result)
		}
	})

	t.Run("tag dry run prints command", func(t *testing.T) {
		tempDir := chdirTemp(t)
		versionFile := filepath.Join(tempDir, "version", "version.go")
		writeGoVersionFile(t, versionFile, "v1.2.3")

		out := captureStdout(t, func() {
			if err := Exec([]string{"gover", "tag", "--dry-run", "--file", versionFile, "--command", "printf tag-{{ .Version }}"}); err != nil {
				t.Fatal(err)
			}
		})

		if strings.TrimSpace(out) != "printf tag-v1.2.3" {
			t.Fatalf("got %q, want rendered tag command", out)
		}
	})

	t.Run("tag dry run json renders command", func(t *testing.T) {
		tempDir := chdirTemp(t)
		versionFile := filepath.Join(tempDir, "version", "version.go")
		writeGoVersionFile(t, versionFile, "v1.2.3")

		out := captureStdout(t, func() {
			if err := Exec([]string{"gover", "tag", "--dry-run", "--json", "--file", versionFile, "--command", "printf tag-{{ .Version }}"}); err != nil {
				t.Fatal(err)
			}
		})

		var result commandResult
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatal(err)
		}
		if !result.DryRun || result.Command != "printf tag-v1.2.3" {
			t.Fatalf("got %+v, want dry-run tag command", result)
		}
	})

	t.Run("release executes bump commit and tag", func(t *testing.T) {
		tempDir := chdirTemp(t)
		versionFile := filepath.Join(tempDir, "version", "version.go")
		writeGoVersionFile(t, versionFile, "v1.2.3")

		out := captureStdout(t, func() {
			if err := Exec([]string{"gover", "release", "--patch", "--file", versionFile, "--commit-command", "printf commit-{{ .Version }}", "--tag-command", "printf tag-{{ .Version }}"}); err != nil {
				t.Fatal(err)
			}
		})

		if !strings.Contains(out, "commit-v1.2.4") || !strings.Contains(out, "tag-v1.2.4") {
			t.Fatalf("got %q, want commit and tag command output", out)
		}
		content, err := os.ReadFile(versionFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(content), `const Version = "v1.2.4"`) {
			t.Fatalf("got %q, want release to bump file", string(content))
		}
	})

	t.Run("release dry run json includes push command", func(t *testing.T) {
		tempDir := chdirTemp(t)
		versionFile := filepath.Join(tempDir, "version", "version.go")
		writeGoVersionFile(t, versionFile, "v1.2.3")

		out := captureStdout(t, func() {
			if err := Exec([]string{"gover", "release", "--dry-run", "--json", "--push", "--patch", "--file", versionFile, "--push-command", "printf push"}); err != nil {
				t.Fatal(err)
			}
		})

		var results []releaseResult
		if err := json.Unmarshal([]byte(out), &results); err != nil {
			t.Fatal(err)
		}
		if len(results) != 1 {
			t.Fatalf("got %d results, want 1", len(results))
		}
		if results[0].Push == nil || !results[0].Push.DryRun || results[0].Push.Command != "printf push" {
			t.Fatalf("got %+v, want dry-run push command", results[0].Push)
		}
	})

	t.Run("check recursive json reports valid files", func(t *testing.T) {
		tempDir := chdirTemp(t)
		writeGoVersionFile(t, filepath.Join(tempDir, "services", "api", "version", "version.go"), "v1.2.3")
		if err := os.MkdirAll(filepath.Join(tempDir, "crates", "core"), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(tempDir, "crates", "core", "Cargo.toml"), []byte(`[package]
name = "core"
version = "2.3.4"
`), 0644); err != nil {
			t.Fatal(err)
		}

		out := captureStdout(t, func() {
			if err := Exec([]string{"gover", "check", "--recursive", "--json", tempDir}); err != nil {
				t.Fatal(err)
			}
		})

		var results []checkResult
		if err := json.Unmarshal([]byte(out), &results); err != nil {
			t.Fatal(err)
		}
		if len(results) != 2 {
			t.Fatalf("got %d results, want 2", len(results))
		}
		for _, result := range results {
			if !result.OK || result.Version == "" {
				t.Fatalf("got %+v, want valid check result", result)
			}
		}
	})

	t.Run("check file validates one explicit target", func(t *testing.T) {
		tempDir := chdirTemp(t)
		versionFile := filepath.Join(tempDir, "version", "version.go")
		writeGoVersionFile(t, versionFile, "v1.2.3")

		out := captureStdout(t, func() {
			if err := Exec([]string{"gover", "check", "--file", versionFile}); err != nil {
				t.Fatal(err)
			}
		})

		if !strings.Contains(out, "v1.2.3") || !strings.Contains(out, versionFile) {
			t.Fatalf("got %q, want explicit check output", out)
		}
	})

	t.Run("check path target validates package directory", func(t *testing.T) {
		tempDir := chdirTemp(t)
		packageDir := filepath.Join(tempDir, "apps", "web")
		if err := os.MkdirAll(packageDir, 0755); err != nil {
			t.Fatal(err)
		}
		packageFile := filepath.Join(packageDir, "package.json")
		if err := os.WriteFile(packageFile, []byte(`{
  "name": "web",
  "version": "1.2.3"
}`), 0644); err != nil {
			t.Fatal(err)
		}

		out := captureStdout(t, func() {
			if err := Exec([]string{"gover", "--lang", "js", "check", packageDir}); err != nil {
				t.Fatal(err)
			}
		})

		if !strings.Contains(out, packageFile) || !strings.Contains(out, "1.2.3") {
			t.Fatalf("got %q, want path-targeted check output", out)
		}
	})

	t.Run("check rejects file with recursive", func(t *testing.T) {
		tempDir := chdirTemp(t)
		versionFile := filepath.Join(tempDir, "version", "version.go")
		writeGoVersionFile(t, versionFile, "v1.2.3")

		err := Exec([]string{"gover", "check", "--recursive", "--file", versionFile})
		if err == nil {
			t.Fatal("got nil, want invalid check flag combination")
		}
		if !strings.Contains(err.Error(), "--file cannot be used") {
			t.Fatalf("got %q, want --file combination error", err)
		}
	})

	t.Run("recursive bump reads package config", func(t *testing.T) {
		tempDir := chdirTemp(t)
		apiFile := filepath.Join(tempDir, "services", "api", "version", "version.go")
		workerFile := filepath.Join(tempDir, "jobs", "worker", "version", "version.go")
		writeGoVersionFile(t, apiFile, "v1.2.3")
		writeGoVersionFile(t, workerFile, "v2.3.4")

		if err := os.WriteFile(filepath.Join(tempDir, "services", "api", ".gover"), []byte(`GOVER_PACKAGE_NAME=api
GOVER_LOCAL_VERSION=true
`), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(tempDir, "jobs", "worker", ".gover"), []byte(`GOVER_PACKAGE_NAME=worker
`), 0644); err != nil {
			t.Fatal(err)
		}

		if err := Exec([]string{"gover", "bump", "--recursive", "--patch", tempDir}); err != nil {
			t.Fatal(err)
		}

		apiOut, err := os.ReadFile(apiFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(apiOut), "package api") || !strings.Contains(string(apiOut), `const version = "v1.2.4"`) {
			t.Fatalf("got %q, want api package config applied", string(apiOut))
		}

		workerOut, err := os.ReadFile(workerFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(workerOut), "package worker") || !strings.Contains(string(workerOut), `const Version = "v2.3.5"`) {
			t.Fatalf("got %q, want worker package config applied", string(workerOut))
		}
	})

	t.Run("recursive bump command flag overrides package config", func(t *testing.T) {
		tempDir := chdirTemp(t)
		versionFile := filepath.Join(tempDir, "services", "api", "version", "version.go")
		writeGoVersionFile(t, versionFile, "v1.2.3")

		if err := os.WriteFile(filepath.Join(tempDir, "services", "api", ".gover"), []byte(`GOVER_PACKAGE_NAME=api
`), 0644); err != nil {
			t.Fatal(err)
		}

		if err := Exec([]string{"gover", "bump", "--recursive", "--patch", "--package", "override", tempDir}); err != nil {
			t.Fatal(err)
		}

		out, err := os.ReadFile(versionFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(out), "package override") {
			t.Fatalf("got %q, want command flag to override package config", string(out))
		}
	})

	t.Run("recursive bump discovers configured custom version file", func(t *testing.T) {
		tempDir := chdirTemp(t)
		versionFile := filepath.Join(tempDir, "apps", "api", "internal", "ver.go")
		writeGoVersionFile(t, versionFile, "v1.2.3")

		if err := os.WriteFile(filepath.Join(tempDir, "apps", "api", ".gover"), []byte(`GOVER_LANG=go
GOVER_VERSION_FILE=internal/ver.go
GOVER_PACKAGE_NAME=internal
`), 0644); err != nil {
			t.Fatal(err)
		}

		if err := Exec([]string{"gover", "bump", "--recursive", "--patch", tempDir}); err != nil {
			t.Fatal(err)
		}

		out, err := os.ReadFile(versionFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(out), "package internal") || !strings.Contains(string(out), `const Version = "v1.2.4"`) {
			t.Fatalf("got %q, want configured custom Go file bumped", string(out))
		}
	})

	t.Run("recursive release reads package command templates", func(t *testing.T) {
		tempDir := chdirTemp(t)
		apiDir := filepath.Join(tempDir, "apps", "api")
		webDir := filepath.Join(tempDir, "apps", "web")
		for _, dir := range []string{apiDir, webDir} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatal(err)
			}
		}
		if err := os.WriteFile(filepath.Join(apiDir, "package.json"), []byte(`{"version":"1.2.3"}`), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(webDir, "package.json"), []byte(`{"version":"2.3.4"}`), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(apiDir, ".gover"), []byte(`COMMIT_COMMAND=printf api-commit-{{ .Version }}
GOVER_TAG_COMMAND=printf api-tag-{{ .Version }}
`), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(webDir, ".gover"), []byte(`COMMIT_COMMAND=printf web-commit-{{ .Version }}
GOVER_TAG_COMMAND=printf web-tag-{{ .Version }}
`), 0644); err != nil {
			t.Fatal(err)
		}

		out := captureStdout(t, func() {
			if err := Exec([]string{"gover", "release", "--recursive", "--dry-run", "--patch", tempDir}); err != nil {
				t.Fatal(err)
			}
		})

		for _, want := range []string{"api-commit-1.2.4", "api-tag-1.2.4", "web-commit-2.3.5", "web-tag-2.3.5"} {
			if !strings.Contains(out, want) {
				t.Fatalf("got %q, want %q", out, want)
			}
		}
	})
}

func chdirTemp(t *testing.T) string {
	t.Helper()

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
	return tempDir
}

func writeGoVersionFile(t *testing.T, file, version string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte(`package version

const Version = "`+version+`"
`), 0644); err != nil {
		t.Fatal(err)
	}
}

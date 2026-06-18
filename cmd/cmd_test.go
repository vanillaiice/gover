package cmd

import (
	"fmt"
	"os"
	"path/filepath"
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
}

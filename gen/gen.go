package gen

import (
	"fmt"
	"os"
)

// perm is the file permission used for writing.
const perm = 0644

// VersionFile generates a file containing the package version.
func VersionFile(packageName, version, outputFile string) (err error) {
	content := fmt.Sprintf(`package %s

// version is the current version of the package.
const version = %q`, packageName, version)

	return os.WriteFile(outputFile, []byte(content), perm)
}

package gen

import (
	"fmt"
	"os"
	"regexp"
)

// PackageJsonData is the package.json data.
var versionRe = regexp.MustCompile(`"version"\s*:\s*"[^"]*"`)

// UpdatePackageVersion updates package.json with the new version.
func UpdatePackageVersion(filePath string, version string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return []byte{}, fmt.Errorf("error reading file: %w", err)
	}

	if !versionRe.Match(content) {
		return []byte{}, fmt.Errorf("no version field found in %s", filePath)
	}

	updatedContent := versionRe.ReplaceAll(content, fmt.Appendf([]byte{}, `"version": "%s"`, version))

	return updatedContent, nil
}

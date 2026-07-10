package gen

import (
	"fmt"
	"os"
	"strconv"

	load "github.com/vanillaiice/gover/v3/load/rust"
)

func UpdateCargoVersion(filePath string, version string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	_, valueStart, valueEnd, err := load.PackageVersionRangeForUpdate(content)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", filePath, err)
	}

	quotedVersion := strconv.Quote(version)
	updatedContent := make([]byte, 0, len(content)-valueEnd+valueStart+len(quotedVersion))
	updatedContent = append(updatedContent, content[:valueStart]...)
	updatedContent = append(updatedContent, quotedVersion...)
	updatedContent = append(updatedContent, content[valueEnd:]...)

	return updatedContent, nil
}

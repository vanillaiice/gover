package load

import (
	"fmt"
	"os"
	"regexp"
)

// VersionData is the version data.
type VersionData struct {
	Version string `json:"version"`
}

// FromFile loads the version data from a go file.
func FromFile(file string) (*VersionData, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`(version|Version)\s*=\s*"([^"]*)"`)
	matches := re.FindStringSubmatch(string(content))
	if len(matches) < 2 {
		return nil, fmt.Errorf("could not find version in %s", file)
	}

	versionData := VersionData{
		Version: matches[2],
	}

	return &versionData, nil
}

// FromFilePanic is the same as FromFile but panics on error.
func FromFilePanic(file string) (vd *VersionData) {
	var err error
	vd, err = FromFile(file)
	if err != nil {
		panic(err)
	}

	return
}

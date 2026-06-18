package load

import (
	"fmt"
	"os"
	"regexp"
)

// FromFile loads the version from a go file.
func FromFile(file string) (string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`(version|Version)\s*=\s*"([^"]*)"`)
	matches := re.FindStringSubmatch(string(content))
	if len(matches) < 2 {
		return "", fmt.Errorf("could not find version in %s", file)
	}

	versionData := matches[2]

	return versionData, nil
}

// FromFilePanic is the same as FromFile but panics on error.
func FromFilePanic(file string) string {
	v, err := FromFile(file)
	if err != nil {
		panic(err)
	}
	return v
}

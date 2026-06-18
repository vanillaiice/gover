package load

import (
	"encoding/json"
	"os"
)

// PackageData represents the package.json data.
type PackageData struct {
	Version string `json:"version"`
}

// FromFile loads the version from file.
func FromFile(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	var versionData PackageData
	if err = json.Unmarshal(data, &versionData); err != nil {
		return "", err
	}

	return versionData.Version, nil
}

// FromFilePanic is the same as FromFile but panics on error.
func FromFilePanic(file string) string {
	v, err := FromFile(file)
	if err != nil {
		panic(err)
	}
	return v
}

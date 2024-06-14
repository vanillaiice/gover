package load

import (
	"encoding/json"
	"os"
)

// VersionData is the version data.
type VersionData struct {
	Version string `json:"version"`
}

// FromFile loads the version data from file.
func FromFile(file string) (*VersionData, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var versionData VersionData
	if err = json.Unmarshal(data, &versionData); err != nil {
		return nil, err
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
